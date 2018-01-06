package httpclient

import (
	"io/ioutil"
	"net/http"
	"net/url"
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
	defer r.Response.Body.Close()
	return ioutil.ReadAll(r.Response.Body)
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
	return r.Request.URL
}
