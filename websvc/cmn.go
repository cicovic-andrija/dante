package websvc

import (
	"fmt"
	"net/http"

	"github.com/cicovic-andrija/dante/atlas"
)

// Client-facing messages and message formats in API response objects.
const (
	CFEmptyTargetInRequest       = "Target cannot be empty string."
	CFEndpointNotFound           = "Endpoint not found."
	CFInternalServerErrorFmt     = "Request %s %s failed because of an internal server error."
	CFInvalidNumberOfProbes      = "Number of requested probes must be a positive integer."
	CFInvalidOperationFmt        = "Operation %s is invalid."
	CFInvalidProbeRequestTypeFmt = "Probe request type must be one of: %s"
	CFMethodNotAllowedFmt        = "Method %s is not allowed."
	CFReqDecodingFailed          = "Failed to decode request body."
	CFCreationFailedFmt          = "Measurement %s creation failed because of a system error."
	CFStatusSuccess              = "Success."
	CFStatusAccepted             = "Accepted."
	CFStatusFailed               = "Failed."
	CFStatusScheduled            = "Scheduled."
	CFStatusOngoing              = "Ongoing."
	CFStatusCompleted            = "Completed."
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
