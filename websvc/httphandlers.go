package websvc

import (
	"fmt"
	"net/http"
)

func (s *server) creditsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "heeeey")
}

func (s *server) invalidEndpointHandler(w http.ResponseWriter, r *http.Request) {
	s.httpWriteError(
		w,
		&ErrorResponse{
			Title:       http.StatusText(http.StatusNotFound),
			Code:        http.StatusNotFound,
			Description: CFEndpointNotFound,
		},
	)
}
