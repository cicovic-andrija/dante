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
		// TODO: Validate .Value?
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
			IsPublic:      false,
			IsOneOff:      true,
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

	resp := &atlas.MeasurementReqResponse{}
	if err = s.makeRequest(httpReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
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
		backendIDs: backendIDs,
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
		period:  10 * time.Second,
		log:     s.log,
	}

	s.taskManager.scheduleTask(task, meas)

	return nil
}

// Intended to be run as a timer task, thus the signature.
func (s *server) updateMeasurementResults(args ...interface{}) (status string, failed bool) {
	// convert generic argument to measurement this task is tracking
	meas := args[0].(*measurement)

	for _, backendID := range meas.backendIDs {
		// fetch results
		req, err := atlas.PrepareRequest(
			atlas.MeasurementResultsURL(backendID),
			&atlas.ReqParams{
				Method: http.MethodGet,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err != nil {
			// TODO: Maybe continue instead?
			return timerTaskFailure(err)
		}
		s.log.info(req.URL.String())

		var results atlas.MeasurementResults
		if err = s.makeRequest(req, &results); err != nil {
			// TODO: Maybe continue instead?
			return timerTaskFailure(err)
		}

		s.log.info("ID: %d, Length: %d", backendID, len(results))
		s.log.info("%+v", results)

		for _, result := range results {
			for _, res := range result.Results {
				s.database.WriteHTTPMeasurementResult(
					meas.BucketName,
					backendID,
					res.RT,
					time.Unix(result.Timestamp, 0),
				)
			}
			// TODO: Handle errors.
		}
	}

	return timerTaskSuccess("done.")
}
