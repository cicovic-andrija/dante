package websvc

import (
	"encoding/json"
	"net/http"
)

const (
	OperationShutdown = "shutdown"

	InvalidOperationFmt = "invalid operation: %s"
)

// Control
type Control struct {
	Operation string `json:"operation"`
}

type status struct {
	Status string `json:"status"`
}

func (s *server) controlHandler(w http.ResponseWriter, r *http.Request) {
	ctrl := &Control{}
	if ok := s.decodeReqBody(w, r, ctrl); !ok {
		return
	}

	switch ctrl.Operation {
	case OperationShutdown:
		s.signalShutdown()
	default:
		s.badRequest(w, CFInvalidOperationFmt, ctrl.Operation)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&status{Status: CFSuccess}) // FIXME: Handle error.
}
