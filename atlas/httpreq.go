package atlas

import (
	"fmt"
	"io"
	"net/http"
)

const (
	URLBase = "https://atlas.ripe.net:443/api/v2"

	CreditsEndpoint = URLBase + "/credits"

	AuthorizationHeader = "Authorization"
	AuthorizationFmt    = "Key %s"

	ContentTypeHeader = "Content-Type"
	ContentType       = "application/json"
)

type ReqParams struct {
	Key string
}

func PrepareRequest(method string, url string, body io.Reader, reqParams *ReqParams) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url, nil)
	default:
		return nil, fmt.Errorf("method %q is invalid or not supported", method)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		AuthorizationHeader,
		fmt.Sprintf(AuthorizationFmt, reqParams.Key),
	)

	return req, nil
}
