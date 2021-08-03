package websvc

import (
	"fmt"

	"github.com/cicovic-andrija/dante/db"
)

// Intended to be run as a timer task, thus the signature.
func (s *server) getCredits(args ...interface{} /* unused */) (status string, failed bool) {
	credit, err := s.httpGetCredits()

	if err == nil {
		err = s.database.WriteCreditBalance(credit.CurrentBalance)
	}

	if err == nil {
		return timerTaskSuccess(fmt.Sprintf("remaining credits: %d", credit.CurrentBalance))
	}

	return timerTaskFailure(err)
}

// Intended to be run as a timer task, thus the signature.
func (s *server) probeDatabase(args ...interface{} /* unused */) (status string, failed bool) {
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
