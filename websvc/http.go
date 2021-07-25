package websvc

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func (s *server) httpInit() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/api/operations", s.operationRequest)

	s.httpServer = &http.Server{
		Addr:     cfg.Net.GetAddr(),
		Handler:  s.router,
		ErrorLog: s.log.errorLogger.backend,
	}
}

func (s *server) httpRun() {
	s.httpWg = &sync.WaitGroup{}
	s.httpWg.Add(1)

	go func() {
		s.log.info("[http] starting the server on %q ...", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			s.log.info("[http] server shut down gracefully")
		} else {
			s.log.err("[http] failed to shut down the server: %v", err)
		}
		s.httpWg.Done()
	}()
}

func (s *server) httpClientInit() {
	s.httpClient = &http.Client{
		Timeout: HTTPClientTimeout,
	}
}
