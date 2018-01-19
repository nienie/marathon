package httpclient

import (
	"net/http"
	"net/url"

	httputil "github.com/nienie/marathon/utils/http"
)

//HTTPResponse ...
type HTTPResponse struct {
	*http.Response
	payload []byte
}

//NewHTTPResponse ...
func NewHTTPResponse(resp *http.Response) *HTTPResponse {
	r := &HTTPResponse{}
	r.Response = resp
	return r
}

//GetPayload ...
func (r *HTTPResponse) GetPayload() ([]byte, error) {
	if len(r.payload) > 0 {
		return r.payload, nil
	}
	var err error
	r.payload, err = httputil.DumpResponseBody(r.Response)
	return r.payload, err
}

//HasPayload ...
func (r *HTTPResponse) HasPayload() bool {
	return r.Response.ContentLength > 0
}

//IsSuccess ...
func (r *HTTPResponse) IsSuccess() bool {
	return r.Response.StatusCode/100 == 2
}

//GetHeaders ...
func (r *HTTPResponse) GetHeaders() map[string][]string {
	return r.Response.Header
}

//GetRequestedURI ...
func (r *HTTPResponse) GetRequestedURI() *url.URL {
	return r.Response.Request.URL
}

//GetStatusCode ...
func (r *HTTPResponse) GetStatusCode() int {
	return r.Response.StatusCode
}
