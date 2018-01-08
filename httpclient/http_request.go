package httpclient

import (
	"io"
	"net/http"
	"net/url"

	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/config"
)

//HTTPRequest ...
type HTTPRequest struct {
	*http.Request
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
	return rr, nil
}

//CreateHTTPRequest ...
func CreateHTTPRequest(r *http.Request, clientConfig config.ClientConfig) *HTTPRequest {
	return &HTTPRequest{
		Request:         r,
		loadBalancerKey: clientConfig.GetPropertyAsString(config.LoadBalancerKey, config.DefaultLoadBalancerKey),
	}
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
func (r *HTTPRequest) ReplaceURI(newURI *url.URL) client.Request {
	return r.SetURI(newURI)
}
