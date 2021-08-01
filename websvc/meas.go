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
		Bucket:     fmt.Sprintf(measBucketNameFmt, id),
		URL:        "",
		backendIDs: backendIDs,
	}

	// update server cache
	s.measCache.insert(meas)

	return meas, nil
}

func (s *server) scheduleWorker(meas *measurement) {
	task := &timerTask{
		name:    meas.Id,
		execute: s.updateMeasurementResults,
		period:  5 * time.Minute,
		log:     s.log,
	}

	s.taskManager.scheduleTask(task, meas)
}

// Intended to be run as a timer task, thus the signature.
func (s *server) updateMeasurementResults(args ...interface{}) (status string, failed bool) {
	meas := args[0].(*measurement)
	s.log.info("Working on measurement %s", meas.Id)
	return timerTaskSuccess("Successful...")
}
