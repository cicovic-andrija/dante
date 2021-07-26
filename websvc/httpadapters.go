package websvc

import (
	"net/http"

	"github.com/cicovic-andrija/dante/util"
)

type Adapter func(h http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// adapter generators

func (s *server) allowMethods(methods ...string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if found := util.SearchForString(r.Method, methods...); !found {
				s.badRequest(w, CFMethodNotAllowedFmt, r.Method)
				return
			}
			// call original handler
			h.ServeHTTP(w, r)
		})
	}
}

// adapters

func (s *server) logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.info(
			"[http] request: %s %s",
			r.Method,
			r.URL.Path,
		)
		// call original handler
		h.ServeHTTP(w, r)
	})
}
