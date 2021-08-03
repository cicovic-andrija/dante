package websvc

import (
	"net/http"

	"github.com/cicovic-andrija/dante/atlas"
	"github.com/cicovic-andrija/dante/db"
)

func (s *server) measurementsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// HTTP GET
	case http.MethodGet:
		measurements := s.measCache.getAll()
		s.httpWriteResponseObject(w, r, http.StatusOK, measurements)

	// HTTP PUT
	case http.MethodPut:
		var (
			err        error
			measReq    = &measurementReq{}
			backendIDs *atlas.MeasurementReqResponse
			meas       *measurement
		)

		if ok := s.decodeReqBody(w, r, measReq); !ok {
			return
		}

		if ok, errMsg := s.validateMeasurementReq(measReq); !ok {
			s.badRequest(w, r, errMsg)
			return
		}

		if backendIDs, err = s.createBackendMeasurements(measReq); err != nil {
			s.internalServerError(w, r, err)
			return
		}

		if meas, err = s.mintMeasurement(backendIDs.Measurements); err != nil {
			s.deleteBackendMeasurements(backendIDs.Measurements...)
			s.internalServerError(w, r, err)
			return
		}

		if err = s.scheduleWorker(meas); err != nil {
			// update server cache
			s.measCache.del(meas.Id)

			s.deleteBackendMeasurements(backendIDs.Measurements...)
			s.internalServerError(w, r, err)
			return
		}

		s.httpWriteResponseObject(w, r, http.StatusCreated, meas)
	}
}

func (s *server) singleMeasurementHandler(w http.ResponseWriter, r *http.Request, routeVars map[string]string) {
	switch r.Method {
	// HTTP GET
	case http.MethodGet:
		meas, ok := s.measCache.get(routeVars["id"])
		if ok {
			s.httpWriteResponseObject(w, r, http.StatusOK, meas)
		} else {
			s.httpWriteResponseObject(w, r, http.StatusNotFound, NotFound)
		}

	// HTTP DELETE
	case http.MethodDelete:
		// TODO
	}
}

func (s *server) creditsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		creditResp = &creditResp{}
	)

	credit, err := s.httpGetCredits()
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	creditResp.CurrentBalance = credit.CurrentBalance
	creditResp.URL = s.database.DataExplorerURL(db.OperationalDataBucket)

	s.httpWriteResponseObject(w, r, http.StatusOK, creditResp)
}

func (s *server) invalidEndpointHandler(w http.ResponseWriter, r *http.Request) {
	s.httpWriteResponseObject(w, r, http.StatusNotFound, NotFound)
}
