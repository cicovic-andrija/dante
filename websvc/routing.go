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

	router.Handle(
		"/api/measurements",
		Adapt(
			http.HandlerFunc(s.measurementsHandler),
			s.logRequest,
			s.allowMethods(http.MethodGet, http.MethodPut),
		),
	)

	router.Handle(
		"/api/measurements/{id:[0-9a-f]+}",
		Adapt(
			variableRouteHandler(s.singleMeasurementHandler),
			s.logRequest,
			s.allowMethods(http.MethodGet, http.MethodDelete),
		),
	)

	router.Handle(
		"/api/measurements/{id:[0-9a-f]+}/control",
		Adapt(
			variableRouteHandler(s.measurementControlHandler),
			s.logRequest,
			s.allowMethods(http.MethodPost),
		),
	)

	// catch all
	router.PathPrefix("/").HandlerFunc(s.invalidEndpointHandler)

	return router
}

func variableRouteHandler(handler func(http.ResponseWriter, *http.Request, map[string]string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, mux.Vars(r))
	}
}
