package websvc

import (
	"fmt"
	"net/http"

	"github.com/cicovic-andrija/dante/atlas"
	"github.com/cicovic-andrija/dante/db"
)

func (s *server) getCredits() (status string, failed bool) {
	var (
		reqParams = &atlas.ReqParams{Key: cfg.Atlas.Auth.Key}
		credit    = &atlas.Credit{}
		req       *http.Request
		err       error
	)

	req, err = atlas.PrepareRequest(http.MethodGet, atlas.CreditsEndpoint, nil, reqParams)

	if err == nil {
		err = s.makeRequest(req, credit)
	}

	if err == nil && credit.Error.Status != 0 {
		err = fmt.Errorf("request failed (%s %d): %s", credit.Error.Title, credit.Error.Status, credit.Error.Detail)
	}

	if err == nil {
		err = s.database.WriteCreditBalance(credit.CurrentBalance)
	}

	if err == nil {
		return timerTaskSuccess(fmt.Sprintf("remaining credits: %d", credit.CurrentBalance))
	}

	return timerTaskFailure(err)
}

func (s *server) probeDatabase() (status string, failed bool) {
	var (
		report = &db.HealthReport{}
		err    error
	)

	err = s.httpGet(s.database.HealthEndpoint(), report)

	if err == nil {
		if report.Status == "pass" {
			return timerTaskSuccess(
				fmt.Sprintf("database is healthy and available on %s: %s",
					s.database.ServerURL(), report.Message),
			)
		}
		err = fmt.Errorf("database is unhealthy with status %q: %s",
			report.Status, report.Message)
	}

	return timerTaskFailure(err)
}
