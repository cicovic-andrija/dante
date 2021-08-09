package websvc

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cicovic-andrija/dante/atlas"
)

// Client-facing messages and message formats in API response objects.
const (
	CFCreationFailedFmt          = "Measurement %s creation failed because of a system error."
	CFEmptyTargetInRequest       = "Target cannot be empty string."
	CFEndpointNotFound           = "Endpoint not found."
	CFEndTimeBeforeStartTime     = "End time cannot be a value earlier than start time."
	CFEndTimeNotSpecified        = "End time not specified."
	CFInternalServerErrorFmt     = "Request %s %s failed because of an internal server error."
	CFIntervalValueTooLarge      = "Interval value too large for the specified time window."
	CFInvalidIntervalValue       = "Interval value not specified or invalid. Value must be a positive integer."
	CFInvalidNumberOfProbes      = "Number of requested probes must be a positive integer."
	CFInvalidOperationFmt        = "Operation %s is invalid."
	CFInvalidProbeRequestTypeFmt = "Probe request type must be one of: %s"
	CFInvalidTimeValueFmt        = "Failed to parse time value: %s."
	CFMeasurementNoStop          = "This measurement cannot be stopped."
	CFMethodNotAllowedFmt        = "Method %s is not allowed."
	CFProbeRequestNotSpecified   = "At least one probe request must be specified."
	CFReqDecodingFailed          = "Failed to decode request body."
	CFResourceNotFound           = "Resource not found."
	CFStartTimeNotSpecified      = "Start time not specified."
	CFTargetNotSpecified         = "At least one target must be specified."

	CFStatusSuccess = "Success."

	CFStatusFailed    = "Failed."
	CFStatusOngoing   = "Ongoing."
	CFStatusQueued    = "Queued."
	CFStatusScheduled = "Scheduled."
	CFStatusStopped   = "Stopped."
)

// This request is issued in multiple places,
// thus putting common code in a separate method.
func (s *server) httpGetCredits() (*atlas.Credit, error) {
	var (
		reqParams = &atlas.ReqParams{Method: http.MethodGet, Key: cfg.Atlas.Auth.Key}
		credit    = &atlas.Credit{}
		req       *http.Request
		err       error
	)

	if req, err = atlas.PrepareRequest(atlas.CreditsEndpoint, reqParams); err != nil {
		return nil, err
	}

	if err = s.makeRequest(req, credit); err != nil {
		return nil, err
	}

	if credit.Error != nil {
		err = fmt.Errorf(
			"client request failed (%s %d): %s",
			credit.Error.Title, credit.Error.Status, credit.Error.Detail,
		)
		credit = nil
	}

	return credit, err
}

// help functions

func backendIDsToStr(backendMeasurements []*backendMeasurement) string {
	strs := make([]string, 0, len(backendMeasurements))
	for _, bm := range backendMeasurements {
		strs = append(strs, strconv.FormatInt(bm.ID, 10))
	}
	return strings.Join(strs, ";")
}

func strToBackendIDs(str string) []int64 {
	strs := strings.Split(str, ";")
	ids := make([]int64, 0, len(strs))
	for _, s := range strs {
		i, _ := strconv.ParseInt(s, 10, 64)
		ids = append(ids, i)
	}
	return ids
}
