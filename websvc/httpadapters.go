package websvc

import (
	"net/http"

	"github.com/cicovic-andrija/dante/util"
)

// Adapter is an HTTP handler that invokes another HTTP handler.
type Adapter func(h http.Handler) http.Handler

// Adapt returnes an HTTP handler enhanced by a number of specified adapters.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func (s *server) allowMethods(methods ...string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if found := util.SearchForString(r.Method, methods...); !found {
				s.badRequest(w, r, CFMethodNotAllowedFmt, r.Method)
				return
			}
			// call original handler
			h.ServeHTTP(w, r)
		})
	}
}

func (s *server) logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.info("%s: accepted", httpReqInfoPrefix(r))

		// call original handler
		h.ServeHTTP(w, r)
	})
}
