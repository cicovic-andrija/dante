package websvc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cicovic-andrija/dante/atlas"
	"github.com/cicovic-andrija/dante/db"
	"github.com/cicovic-andrija/dante/util"
)

// TODO: Stop measurement.
// TODO: Stop measurement worker.
// TODO: Restore.
// TODO: Limit lines to 80 or 120 characters.

const (
	measDefDescrFmt   = "IPv4/HTTP measurement for target %s."
	measBucketNameFmt = "meas-%s"
	measIDHexLength   = 9
)

func freshMeasurementID() (id string, err error) {
	id, err = util.RandHexString(measIDHexLength)
	return
}

func (s *server) validateMeasurementReq(req *measurementReq) (bool, string) {
	var (
		startTime time.Time
		endTime   time.Time
		err       error
	)

	if len(req.Targets) == 0 {
		return false, CFTargetNotSpecified
	}

	for _, target := range req.Targets {
		if target == "" {
			return false, CFEmptyTargetInRequest
		}
	}

	if len(req.ProbeRequests) == 0 {
		return false, CFProbeRequestNotSpecified
	}

	for _, probeReq := range req.ProbeRequests {
		if probeReq.Requested < 1 {
			return false, CFInvalidNumberOfProbes
		}
		if found := util.SearchForString(probeReq.Type, atlas.ValidProbeRequestTypes...); !found {
			return false, fmt.Sprintf(CFInvalidProbeRequestTypeFmt, atlas.ValidProbeRequestTypesStr)
		}
	}

	if req.StartTimeRFC3339 == "" {
		return false, CFStartTimeNotSpecified
	}

	if startTime, err = time.Parse(time.RFC3339, req.StartTimeRFC3339); err != nil {
		return false, fmt.Sprintf(CFInvalidTimeValueFmt, req.StartTimeRFC3339)
	}

	if req.StopTimeRFC3339 == "" {
		return false, CFEndTimeNotSpecified
	}

	if endTime, err = time.Parse(time.RFC3339, req.StopTimeRFC3339); err != nil {
		return false, fmt.Sprintf(CFInvalidTimeValueFmt, req.StopTimeRFC3339)
	}

	if endTime.Before(startTime) {
		return false, CFEndTimeBeforeStartTime
	}

	req.startTimeUnix = startTime.Unix()
	req.stopTimeUnix = endTime.Unix()

	if req.IntervalSec <= 0 {
		return false, CFInvalidIntervalValue
	}

	if req.IntervalSec >= (req.stopTimeUnix - req.startTimeUnix) {
		return false, CFIntervalValueTooLarge
	}

	return true, ""
}

func (s *server) measurementCreationWorkflow(req *measurementReq, id string) {
	var (
		resp *atlas.MeasurementReqResponse
		meas *measurement
		err  error
	)

	commitFailedMeasurement := func() {
		s.measCache.insert(
			&measurement{
				ID:     id,
				Status: CFStatusFailed,
				Reason: fmt.Sprintf(CFCreationFailedFmt, id),
			},
		)
	}

	if resp, err = s.createBackendMeasurements(req); err != nil {
		s.log.err("[mgmt %s] backend measurement creation failed: %v", id, err)
		commitFailedMeasurement()
		return
	}

	if meas, err = s.mintMeasurement(resp.Measurements, req.Description, id); err != nil {
		s.stopBackendMeasurements(resp.Measurements...)
		s.log.err("[mgmt %s] internal creation failed: %v", id, err)
		commitFailedMeasurement()
		return
	}

	if err = s.scheduleWorker(meas); err != nil {
		s.cleanupMeasurement(meas)
		s.stopBackendMeasurements(resp.Measurements...)
		s.log.err("[mgmt %s] failed to schedule worker: %v", id, err)
		commitFailedMeasurement()
		return
	}

	// commit successfully created measurement
	s.measCache.insert(meas)
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
			IsOneOff:      false,
			StartTime:     req.startTimeUnix,
			StopTime:      req.stopTimeUnix,
			Interval:      req.IntervalSec,
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

	if resp.Error != nil {
		return nil, fmt.Errorf(
			"client request failed (%s %d): %s",
			resp.Error.Title, resp.Error.Status, resp.Error.Detail,
		)
	}

	return resp, nil
}

// Best effort, ignore all errors.
func (s *server) stopBackendMeasurements(backendIDs ...int64) {
	for _, id := range backendIDs {
		req, err := atlas.PrepareRequest(
			atlas.MeasurementURL(id),
			&atlas.ReqParams{
				Method: http.MethodDelete,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err == nil {
			s.makeRequest(req, nil)
		}
	}
}

func (s *server) mintMeasurement(backendIDs []int64, description string, id string) (*measurement, error) {
	var (
		err error
	)

	bucketName := fmt.Sprintf(measBucketNameFmt, id)

	meas := &measurement{
		ID:          id,
		BucketName:  bucketName,
		Description: description,
		Status:      CFStatusQueued,
		URL:         s.database.DataExplorerURL(bucketName),
	}

	// issue requests to the backend API for
	// backend measurement details
	if err = s.fetchRetainBackendDetails(meas, backendIDs); err != nil {
		return nil, err
	}

	return meas, nil
}

func (s *server) fetchRetainBackendDetails(meas *measurement, backendIDs []int64) error {
	var (
		req *http.Request
		err error
	)

	meas.BackendMeasurements = make([]*backendMeasurement, 0, len(backendIDs))

	for _, id := range backendIDs {
		req, err = atlas.PrepareRequest(
			atlas.MeasurementURL(id),
			&atlas.ReqParams{
				Method: http.MethodGet,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err != nil {
			return err
		}

		resp := &atlas.Measurement{}
		if err = s.makeRequest(req, resp); err != nil {
			return err
		}

		bm := &backendMeasurement{
			ID:       id,
			Target:   resp.Target,
			TargetIP: resp.TargetIP,

			stopped: resp.Status.ID > atlas.MeasurementStatusIDOngoing,
		}
		meas.BackendMeasurements = append(meas.BackendMeasurements, bm)
	}

	return nil
}

func (s *server) cleanupMeasurement(meas *measurement) {
	// delete bucket if it exists
	if meas.bucket != nil {
		err := s.database.DeleteBucket(meas.bucket)
		if err != nil {
			s.log.err("[mgmt %s] failed to delete bucket %s", meas.ID, meas.BucketName)
		}
		meas.bucket = nil
	}
}

func (s *server) scheduleWorker(meas *measurement) error {
	// first ensure there is a bucket for writing data
	if meas.bucket == nil {
		if bck, err := s.database.EnsureBucket(meas.BucketName); err == nil {
			meas.bucket = bck
		} else {
			return err
		}
	}

	task := &timerTask{
		name:    meas.ID,
		execute: s.updateMeasurementResults,
		period:  5 * time.Minute,
		log:     s.log,
	}

	meas.Status = CFStatusScheduled

	// after this point, the worker task owns the pointer to meas
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

	for _, backend := range meas.BackendMeasurements {
		// ignore backend measurement if flagged as stopped in a past iteration
		if backend.stopped {
			continue
		}

		req, err := atlas.PrepareRequest(
			atlas.MeasurementResultsURL(backend.ID),
			&atlas.ReqParams{
				Method: http.MethodGet,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err != nil {
			recordError(fmt.Errorf("prepare request failed for %d: %v", backend.ID, err))
			continue
		}

		var results atlas.MeasurementResults
		if err = s.makeRequest(req, &results); err != nil {
			recordError(fmt.Errorf("request failed for %d: %v", backend.ID, err))
			continue
		}

		for _, probeResults := range results {
			if err = s.processProbeResults(&probeResults, backend, meas.BucketName); err != nil {
				recordError(err)
				continue
			}
		}

		// fetch backend measurement status from the API and update internal state
		req, err = atlas.PrepareRequest(
			atlas.MeasurementURL(backend.ID),
			&atlas.ReqParams{
				Method: http.MethodGet,
				Key:    cfg.Atlas.Auth.Key,
			},
		)
		if err != nil {
			recordError(fmt.Errorf("prepare request for metadata failed for %d: %v", backend.ID, err))
			continue
		}

		resp := &atlas.Measurement{}
		if err = s.makeRequest(req, resp); err != nil {
			recordError(fmt.Errorf("request for metadata failed for %d: %v", backend.ID, err))
			continue
		}

		if resp.Status.ID > atlas.MeasurementStatusIDOngoing {
			// this should happen ONLY in case of no errors in this loop iteration
			backend.stopped = true
		} else if resp.Status.ID == atlas.MeasurementStatusIDOngoing {
			// other threads may update measurement status,
			// so do it in a locked code section
			s.measCache.Lock()
			if meas.Status == CFStatusScheduled {
				meas.Status = CFStatusOngoing
			}
			s.measCache.Unlock()
		}
	}

	stopWorker := true
	for _, backend := range meas.BackendMeasurements {
		if stopWorker = stopWorker && backend.stopped; !stopWorker {
			break
		}
	}
	if stopWorker {
		// other threads may update measurement status,
		// so do it in a locked code section
		s.measCache.Lock()
		// stop worker only if it's not stopped externally (by another thread)
		if meas.Status == CFStatusScheduled || meas.Status == CFStatusOngoing {
			s.taskManager.stopTask(meas.ID)
			meas.Status = CFStatusCompleted
		}
		s.measCache.Unlock()
	}

	if accuError == nil {
		return timerTaskSuccess("no errors")
	} else {
		return timerTaskFailure(accuError)
	}
}

// TODO: Implement smart updating.
func (s *server) processProbeResults(probeResults *atlas.ProbeMeasurementResults, backend *backendMeasurement, bucketName string) error {
	var (
		probe *atlas.Probe
		err   error
	)

	if probe, err = s.getProbe(probeResults.ProbeID); err != nil {
		return fmt.Errorf("probe info request failed for probe %d and measurement %d: %v", probeResults.ProbeID, backend.ID, err)
	}

	for _, result := range probeResults.Results {
		httpData := &db.HTTPData{
			BackendID:     backend.ID,
			ProbeID:       probe.ID,
			ASN:           probe.ASNv4,
			Country:       probe.CountryCode,
			Target:        backend.Target,
			TargetIP:      backend.TargetIP,
			RoundTripTime: result.RT,
			BodySize:      result.BodySize,
			HeaderSize:    result.HeaderSize,
			StatusCode:    result.Result,
			Timestamp:     time.Unix(probeResults.Timestamp, 0),
		}
		if err = s.database.WriteHTTPMeasurementResult(bucketName, httpData); err != nil {
			// do not continue, assume others will fail too
			return fmt.Errorf("writing data point failed for %d: %v", backend.ID, err)
		}
	}

	return nil
}

func (s *server) stopMeasurement(id string) (stat *status, err error) {
	// stop the worker task in a locked code section because
	// the measurement object is being updated by the non-worker tread
	s.measCache.Lock()
	meas := s.measCache.measurements[id]

	if meas.Status != CFStatusScheduled && meas.Status != CFStatusOngoing {
		s.measCache.Unlock()
		return &status{Status: CFStatusFailed, Explanation: CFMeasurementNoStop}, nil
	}

	s.taskManager.stopTask(id)
	meas.Status = CFStatusStopped
	s.measCache.Unlock()

	// best effort, ignore errors
	for _, backendMeasurement := range meas.BackendMeasurements {
		s.stopBackendMeasurements(backendMeasurement.ID)
	}

	return &status{Status: CFStatusSuccess}, nil
}

func (s *server) getProbe(id int64) (*atlas.Probe, error) {
	var (
		probe *atlas.Probe
		req   *http.Request
		ok    bool
		err   error
	)

	if probe, ok = s.probeInfo.lookup(id); ok {
		return probe, nil
	}

	req, err = atlas.PrepareRequest(
		atlas.ProbeURL(id),
		&atlas.ReqParams{
			Method: http.MethodGet,
			Key:    cfg.Atlas.Auth.Key,
		},
	)
	if err != nil {
		return nil, err
	}

	probe = &atlas.Probe{}
	if err = s.makeRequest(req, probe); err != nil {
		return nil, err
	}

	// update probe cache
	s.probeInfo.insert(probe)
	s.log.info("[mgmt] probe info cached: id=%d country=%s asn=%d",
		probe.ID, probe.CountryCode, probe.ASNv4)

	return probe, nil
}
