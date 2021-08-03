package websvc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cicovic-andrija/dante/atlas"
	"github.com/cicovic-andrija/dante/util"
)

const (
	measDefDescrFmt   = "IPv4/HTTP measurement for target %s."
	measBucketNameFmt = "meas-%s"
	measIDHexLength   = 9
)

func (s *server) validateMeasurementReq(req *measurementReq) (bool, string) {
	for _, target := range req.Targets {
		if target == "" {
			return false, CFEmptyTargetInRequest
		}
	}

	for _, probeReq := range req.ProbeRequests {
		if probeReq.Requested < 1 {
			return false, CFInvalidNumberOfProbes
		}
		if found := util.SearchForString(probeReq.Type, atlas.ValidProbeRequestTypes...); !found {
			return false, fmt.Sprintf(CFInvalidProbeRequestTypeFmt, atlas.ValidProbeRequestTypesStr)
		}
	}

	return true, ""
}

func (s *server) createBackendMeasurements(req *measurementReq) (*atlas.MeasurementReqResponse, error) {
	var (
		backendReq = &atlas.MeasurementRequest{}
		httpReq    *http.Request
		err        error
	)

	for _, target := range req.Targets {
		def := &atlas.MeasurementDefinition{
			Type:          atlas.MeasHTTP,
			AddressFamily: atlas.IPv4,
			Target:        target,
			Description:   fmt.Sprintf(measDefDescrFmt, target),
			Interval:      60,
			IsPublic:      false,
			IsOneOff:      false,
		}
		backendReq.Definitions = append(backendReq.Definitions, def)
	}

	for _, probeReq := range req.ProbeRequests {
		probes := &atlas.ProbeRequest{
			Requested: probeReq.Requested,
			Type:      probeReq.Type,
			Value:     probeReq.Value,
		}
		backendReq.Probes = append(backendReq.Probes, probes)
	}

	httpReq, err = atlas.PrepareRequest(
		atlas.MeasurementsEndpoint,
		&atlas.ReqParams{
			Method: http.MethodPost,
			Key:    cfg.Atlas.Auth.Key,
			Body:   backendReq,
		},
	)

	if err != nil {
		return nil, err
	}

	// TODO: Handle errors.

	resp := &atlas.MeasurementReqResponse{}
	if err = s.makeRequest(httpReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Best effort, ignore all errors.
func (s *server) deleteBackendMeasurements(backendIDs ...int64) {
	for _, id := range backendIDs {
		req, _ := atlas.PrepareRequest(
			atlas.MeasurementURL(id),
			&atlas.ReqParams{
				Method: http.MethodDelete,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		s.makeRequest(req, nil)
	}
}

func (s *server) mintMeasurement(backendIDs []int64) (*measurement, error) {
	var (
		id  string
		err error
	)

	if id, err = util.RandHexString(measIDHexLength); err != nil {
		return nil, err
	}

	meas := &measurement{
		Id:         id,
		BucketName: fmt.Sprintf(measBucketNameFmt, id),
		URL:        "",
		BackendIDs: backendIDs,
	}

	// update server cache
	s.measCache.insert(meas)

	return meas, nil
}

func (s *server) scheduleWorker(meas *measurement) error {
	// first ensure there is a bucket for writing data
	if meas.bucket == nil {
		if bck, err := s.database.EnsureBucket(meas.BucketName); err != nil {
			return err
		} else {
			meas.bucket = bck
		}
	}

	task := &timerTask{
		name:    meas.Id,
		execute: s.updateMeasurementResults,
		period:  1 * time.Minute,
		log:     s.log,
	}

	s.taskManager.scheduleTask(task, meas)

	return nil
}

// Intended to be run as a timer task, thus the signature.
func (s *server) updateMeasurementResults(args ...interface{}) (status string, failed bool) {
	// convert generic argument to measurement this task is tracking
	meas := args[0].(*measurement)

	var accuError error
	recordError := func(err error) {
		if accuError == nil {
			accuError = fmt.Errorf("some errors were encountered: {%v}", err)
		} else {
			accuError = fmt.Errorf("%v;{%v}", accuError, err)
		}
	}

	for _, id := range meas.BackendIDs {
		req, err := atlas.PrepareRequest(
			atlas.MeasurementResultsURL(id),
			&atlas.ReqParams{
				Method: http.MethodGet,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err != nil {
			recordError(fmt.Errorf("prepare request failed for %d: %v", id, err))
			continue
		}

		var results atlas.MeasurementResults
		if err = s.makeRequest(req, &results); err != nil {
			recordError(fmt.Errorf("request failed for %d: %v", id, err))
			continue
		}

		// TODO: Implement smart updating.
		for _, probeResult := range results {
			for _, result := range probeResult.Results {
				err = s.database.WriteHTTPMeasurementResult(
					meas.BucketName,
					id,
					result.RT,
					time.Unix(probeResult.Timestamp, 0),
				)
				if err != nil {
					recordError(fmt.Errorf("writing data point failed for %d: %v", id, err))
					continue
				}
			}
		}
	}

	if accuError == nil {
		return timerTaskSuccess("no errors")
	} else {
		return timerTaskFailure(accuError)
	}
}
