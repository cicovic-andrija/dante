package atlas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ReqParams specifies parameters needed to prepare
// an HTTP request for the Atlas API.
type ReqParams struct {
	Method string
	Key    string
	Body   interface{}
}

// PrepareRequest creates a new HTTP request, serializes the request body,
// and set the headers expected by the Atlas API.
func PrepareRequest(url string, reqParams *ReqParams) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	switch reqParams.Method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
	case http.MethodPost:
		var b []byte
		b, err = json.Marshal(reqParams.Body)
		if err == nil {
			req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		}
	default:
		return nil, fmt.Errorf("method %q is invalid or not supported", reqParams.Method)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		AuthorizationHeader,
		fmt.Sprintf(AuthorizationFmt, reqParams.Key),
	)
	req.Header.Set(
		ContentTypeHeader,
		ContentType,
	)

	return req, nil
}
