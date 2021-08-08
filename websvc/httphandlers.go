package websvc

import (
	"fmt"
	"net/http"

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
			err     error
			measID  string
			measReq = &measurementReq{}
		)

		if ok := s.decodeReqBody(w, r, measReq); !ok {
			return
		}

		if ok, errMsg := s.validateMeasurementReq(measReq); !ok {
			s.badRequest(w, r, errMsg)
			return
		}

		// create ID early because it must be returned in HTTP response
		if measID, err = freshMeasurementID(); err != nil {
			s.internalServerError(w, r, err)
			return
		}

		// create measurement in a dedicated thread
		go s.measurementCreationWorkflow(measReq, measID)

		s.httpWriteResponseObject(
			w, r, http.StatusAccepted,
			&status{Status: CFStatusQueued, ID: measID},
		)
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
			s.httpWriteResponseObject(w, r, http.StatusNotFound, ResourceNotFound)
		}

	// HTTP DELETE
	case http.MethodDelete:
		id := routeVars["id"]
		if _, found := s.measCache.get(id); !found {
			s.httpWriteResponseObject(w, r, http.StatusNotFound, ResourceNotFound)
			return
		}

		status := s.deleteMeasurement(id)
		switch status {
		case http.StatusNotFound:
			s.httpWriteResponseObject(w, r, status, ResourceNotFound)
		case http.StatusNoContent:
			w.WriteHeader(status)
		default:
			s.internalServerError(w, r, fmt.Errorf("unexpected status code: %d", status))
		}
	}
}

func (s *server) measurementControlHandler(w http.ResponseWriter, r *http.Request, routeVars map[string]string) {
	// HTTP POST
	ctrl := &control{}
	if ok := s.decodeReqBody(w, r, ctrl); !ok {
		return
	}

	if ctrl.Operation != OperationStop {
		s.badRequest(w, r, CFInvalidOperationFmt, ctrl.Operation)
		return
	}

	status, respObj := s.stopMeasurement(routeVars["id"])
	s.httpWriteResponseObject(w, r, status, respObj)
}

func (s *server) creditsHandler(w http.ResponseWriter, r *http.Request) {
	// HTTP GET
	creditResp := &creditResp{}

	credit, err := s.httpGetCredits()
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	creditResp.CurrentBalance = credit.CurrentBalance
	creditResp.URL = s.database.DataExplorerURL(db.SystemBucket)

	s.httpWriteResponseObject(w, r, http.StatusOK, creditResp)
}

func (s *server) invalidEndpointHandler(w http.ResponseWriter, r *http.Request) {
	s.httpWriteResponseObject(w, r, http.StatusNotFound, NotFound)
}
