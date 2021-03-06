package httpclient

import (
	"io"
	"net/http"
	"net/url"

	"github.com/nienie/marathon/config"
	httputil "github.com/nienie/marathon/utils/http"
)

//HTTPRequest ...
type HTTPRequest struct {
	*http.Request
	body            []byte
	loadBalancerKey interface{}
}

//NewHTTPRequest ...
func NewHTTPRequest(method, urlStr string, body io.Reader, loadBalancerKey interface{}) (*HTTPRequest, error) {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	rr := &HTTPRequest{
		Request:         r,
		loadBalancerKey: loadBalancerKey,
	}
	if r.Body != nil {
		rr.body, err = httputil.DumpRequestBody(r)
		if err != nil {
			return nil, err
		}
	}
	return rr, nil
}

//CreateHTTPRequest ...
func CreateHTTPRequest(r *http.Request, requestConfig config.ClientConfig) *HTTPRequest {
	rr := &HTTPRequest{
		Request: r,
	}
	if requestConfig != nil {
		rr.loadBalancerKey = requestConfig.GetPropertyAsString(config.LoadBalancerKey, config.DefaultLoadBalancerKey)
	}
	if r.Body != nil {
		rr.body, _ = httputil.DumpRequestBody(r)
	}
	return rr
}

//GetURI ...
func (r *HTTPRequest) GetURI() *url.URL {
	return r.Request.URL
}

//SetURI ...
func (r *HTTPRequest) SetURI(uri *url.URL) *HTTPRequest {
	r.Request.URL = uri
	return r
}

//GetLoadBalancerKey ...
func (r *HTTPRequest) GetLoadBalancerKey() interface{} {
	return r.loadBalancerKey
}

//SetLoadBalancerKey ...
func (r *HTTPRequest) SetLoadBalancerKey(loadBalancerKey interface{}) *HTTPRequest {
	r.loadBalancerKey = loadBalancerKey
	return r
}

//GetRawRequest ...
func (r *HTTPRequest) GetRawRequest() *http.Request {
	return r.Request
}

//ReplaceURI ...
func (r *HTTPRequest) ReplaceURI(newURI *url.URL) {
	r.SetURI(newURI)
}

//GetBodyContents ...
func (r *HTTPRequest) GetBodyContents() []byte {
	return r.body
}

//GetBodyLength ...
func (r *HTTPRequest) GetBodyLength() int {
	return len(r.body)
}

//GetHeaders ...
func (r *HTTPRequest) GetHeaders() map[string][]string {
	return r.Header
}
