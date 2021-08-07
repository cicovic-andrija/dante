package websvc

import (
	"net/http"
)

// Control constants.
const (
	OperationShutdown = "shutdown"

	InvalidOperationFmt = "invalid operation: %s"
)

func (s *server) controlHandler(w http.ResponseWriter, r *http.Request) {
	ctrl := &control{}
	if ok := s.decodeReqBody(w, r, ctrl); !ok {
		return
	}

	switch ctrl.Operation {
	case OperationShutdown:
		s.signalShutdown()
	default:
		s.badRequest(w, r, CFInvalidOperationFmt, ctrl.Operation)
		return
	}

	s.httpWriteResponseObject(w, r, http.StatusOK, &status{Status: CFStatusSuccess})
}
