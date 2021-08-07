package websvc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	httpClientTimeout = 25 * time.Second
)

// ErrorResponse specifies a response object sent in case
// of an error encountered during request processing.
type ErrorResponse struct {
	Title       string `json:"title"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// NotFound is a standard ErrorResponse indicating that
// an endpoint or a page is not found.
var NotFound = &ErrorResponse{
	Title:       http.StatusText(http.StatusNotFound),
	Code:        http.StatusNotFound,
	Description: CFEndpointNotFound,
}

// ResourceNotFound is a standard ErrorResponse indicating that
// a resource is not found.
var ResourceNotFound = &ErrorResponse{
	Title:       http.StatusText(http.StatusNotFound),
	Code:        http.StatusNotFound,
	Description: CFResourceNotFound,
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

	if err == nil && v != nil {
		err = json.NewDecoder(res.Body).Decode(&v)
	}

	return err
}

func (s *server) decodeReqBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	var err error
	if err = json.NewDecoder(r.Body).Decode(v); err != nil {
		s.badRequest(w, r, CFReqDecodingFailed)
	}
	return err == nil
}

func (s *server) badRequest(w http.ResponseWriter, r *http.Request, msgFmt string, v ...interface{}) {
	errResp := &ErrorResponse{
		Title:       http.StatusText(http.StatusBadRequest),
		Code:        http.StatusBadRequest,
		Description: fmt.Sprintf(msgFmt, v...),
	}
	s.httpWriteResponseObject(w, r, http.StatusBadRequest, errResp)
}

func (s *server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.log.err("%s: internal server error: %v", httpReqInfoPrefix(r), err)
	errResp := &ErrorResponse{
		Title:       http.StatusText(http.StatusInternalServerError),
		Code:        http.StatusInternalServerError,
		Description: fmt.Sprintf(CFInternalServerErrorFmt, r.Method, r.URL.String()),
	}
	s.httpWriteResponseObject(w, r, http.StatusInternalServerError, errResp)
}

func (s *server) httpWriteResponseObject(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.log.err("%s: response object encoding failed: %v", httpReqInfoPrefix(r), err)
	}
}

func httpReqInfoPrefix(r *http.Request) string {
	return fmt.Sprintf("[http] request %s --> %s %s ", r.RemoteAddr, r.Method, r.URL.String())
}
