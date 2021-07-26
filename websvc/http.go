package websvc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	httpClientTimeout = 15 * time.Second
)

// ErrorResponse
type ErrorResponse struct {
	Title       string `json:"title"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (s *server) httpInit() {
	s.httpClientInit()
	s.router = s.initRouter()
	s.httpServer = &http.Server{
		Addr:     cfg.Net.GetAddr(),
		Handler:  s.router,
		ErrorLog: s.log.errorLogger.backend,
	}
}

func (s *server) httpClientInit() {
	s.httpClient = &http.Client{
		Timeout: httpClientTimeout,
	}
}

func (s *server) runHTTP() {
	s.httpWg = &sync.WaitGroup{}
	s.httpWg.Add(1)

	go func() {
		s.log.info("[http] starting the server on %s ...", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			s.log.info("[http] server shut down gracefully")
		} else {
			s.log.err("[http] failed to shut down the server: %v", err)
		}
		s.httpWg.Done()
	}()
}

func (s *server) httpGet(endpoint string, v interface{}) error {
	var (
		req *http.Request
		err error
	)

	req, err = http.NewRequest(http.MethodGet, endpoint, nil)

	if err == nil {
		err = s.makeRequest(req, v)
	}

	return err
}

func (s *server) makeRequest(req *http.Request, v interface{}) error {
	var (
		res *http.Response
		err error
	)

	res, err = s.httpClient.Do(req)

	if err == nil {
		err = json.NewDecoder(res.Body).Decode(&v)
	}

	return err
}

func (s *server) decodeReqBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	var err error
	if err = json.NewDecoder(r.Body).Decode(v); err != nil {
		s.badRequest(w, CFReqDecodingFailed)
	}
	return err == nil
}

func (s *server) badRequest(w http.ResponseWriter, msgFmt string, v ...interface{}) {
	errResp := &ErrorResponse{
		Title:       http.StatusText(http.StatusBadRequest),
		Code:        http.StatusBadRequest,
		Description: fmt.Sprintf(msgFmt, v...),
	}
	s.httpWriteError(w, errResp)
}

func (s *server) httpWriteError(w http.ResponseWriter, e *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		s.log.err("[http] error response encoding failed: %v", err)
	}
}
