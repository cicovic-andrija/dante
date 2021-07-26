package websvc

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HTTP endpoint definitions.

func (s *server) initRouter() http.Handler {
	// library-provided router
	router := mux.NewRouter()

	router.Handle(
		"/api/control",
		Adapt(
			http.HandlerFunc(s.controlHandler), // main handler (last executed)
			s.logRequest,                       // last executed adapter
			s.allowMethods(http.MethodPost),    // first executed adapter
		),
	)

	router.Handle(
		"/api/credits",
		Adapt(
			http.HandlerFunc(s.creditsHandler),
			s.logRequest,
			s.allowMethods(http.MethodGet),
		),
	)

	// catch all
	router.PathPrefix("/").HandlerFunc(s.invalidEndpointHandler)

	return router
}
