package httpclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/errors"
	"github.com/nienie/marathon/loadbalancer"
	"github.com/nienie/marathon/logger"

	transport "github.com/mreiferson/go-httpclient"
)

var (
	loggerAfterHook AfterHTTHook = func(ctx context.Context, req *HTTPRequest, resp *HTTPResponse, err error) {
		var (
			requestBody  string
			responseBody string
			format       = "method=%s||host=%s||uri=%s||args=%s||body=%v||request_header=%v||response=%v||status_code=%d||response_header=%v||err=%v"
		)
		if req.GetBodyLength() > 8196 { //large than 8K, perhaps it's a file, so do not log it
			requestBody = "<large request body>"
		} else {
			requestBody = string(req.GetBodyContents())
		}

		if err != nil || resp == nil {
			logger.Warnf(ctx, format, req.Method, req.URL.Host, req.URL.Path, req.URL.RawQuery,
				requestBody, req.Header, nil, 0, nil, err)
			return
		}

		payload, _ := resp.GetPayload()
		if len(payload) > 8196 { //large than 8K, perhaps it's a file, so do not log it
			responseBody = "<large resposne>"
		} else {
			responseBody = string(payload)
		}

		logger.Infof(ctx, format, req.Method, req.URL.Host, req.URL.Path, req.URL.RawQuery, requestBody,
			req.Header, responseBody, resp.StatusCode, resp.Header, err)
		return
	}
)

//BeforeHTTPHook ...
type BeforeHTTPHook func(context.Context, *HTTPRequest)

//AfterHTTHook ...
type AfterHTTHook func(context.Context, *HTTPRequest, *HTTPResponse, error)

//LoadBalancerHTTPClient ...
type LoadBalancerHTTPClient struct {
	*http.Client
	*loadbalancer.BaseLoadBalancerClient
	HTTPClientName string
	Transport      *transport.Transport
	BeforeHooks    []BeforeHTTPHook
	AfterHooks     []AfterHTTHook
}

//NewHTTPLoadBalancerClient ...
func NewHTTPLoadBalancerClient(clientConfig config.ClientConfig, lb loadbalancer.LoadBalancer) *LoadBalancerHTTPClient {
	//create load balancer context
	loadBalancerContext := loadbalancer.NewLoadBalancerContext(clientConfig, lb)
	//create load balancer client
	loadBalancerClient := &loadbalancer.BaseLoadBalancerClient{
		loadBalancerContext,
		nil,
	}
	//create transport
	trans := &transport.Transport{
		ConnectTimeout:   clientConfig.GetPropertyAsDuration(config.ConnectTimeout, config.DefaultConnectTimeout),
		ReadWriteTimeout: clientConfig.GetPropertyAsDuration(config.ReadWriteTimeout, config.DefaultReadWriteTimeout),
		RequestTimeout:   clientConfig.GetPropertyAsDuration(config.RequestTimeout, config.DefaultRequestTimeout),
	}
	//create original http.client
	originalClient := &http.Client{
		Transport: trans,
	}
	//create http client with load balancer
	httpClient := &LoadBalancerHTTPClient{
		Client:                 originalClient,
		BaseLoadBalancerClient: loadBalancerClient,
		HTTPClientName:         clientConfig.GetClientName(),
		Transport:              trans,
		BeforeHooks:            make([]BeforeHTTPHook, 0),
		AfterHooks:             []AfterHTTHook{loggerAfterHook},
	}
	//load balancer context correlate with http client
	loadBalancerClient.Client = httpClient
	return httpClient
}

//Do ...
func (c *LoadBalancerHTTPClient) Do(ctx context.Context, request *HTTPRequest, requestConfig config.ClientConfig) (*http.Response, error) {
	if request == nil || request.Request == nil {
		return nil, fmt.Errorf("wrong type, request is nil")
	}
	c.beforeHTTPHook(ctx, request)
	resp, err := c.BaseLoadBalancerClient.ExecuteWithLoadBalancer(ctx, request, requestConfig)
	if err != nil || resp == nil {
		c.afterHTTPHook(ctx, request, nil, err)
		return nil, err
	}
	response := resp.(*HTTPResponse)
	c.afterHTTPHook(ctx, request, response, err)
	return response.Response, nil
}

//Execute Do not Directly Use...
func (c *LoadBalancerHTTPClient) Execute(ctx context.Context, request client.Request, requestConfig config.ClientConfig) (client.Response, error) {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return nil, errors.NewClientError(errors.General, fmt.Errorf("wrong type, type must be *HTTPRquest, type=%T", request))
	}
	return c.ExecuteHTTP(ctx, req, requestConfig)
}

//ExecuteHTTP Do not Directly Use...
func (c *LoadBalancerHTTPClient) ExecuteHTTP(ctx context.Context, request *HTTPRequest, requestConfig config.ClientConfig) (*HTTPResponse, error) {
	response, err := c.Client.Do(request.GetRawRequest())
	if err != nil {
		return nil, errors.ConvertError(err)
	}
	if response.StatusCode == http.StatusBadGateway ||
		response.StatusCode == http.StatusServiceUnavailable ||
		response.StatusCode == http.StatusGatewayTimeout { //502/503/504
		return nil, errors.NewClientError(errors.ServerThrottled, fmt.Errorf("http status code = %d", response.StatusCode))
	}
	return NewHTTPResponse(response), nil
}

//Shutdown ...
func (c *LoadBalancerHTTPClient) Shutdown() {
	c.Transport.Close()
}

//RegisterBeforeHook ...
func (c *LoadBalancerHTTPClient) RegisterBeforeHook(hooks ...BeforeHTTPHook) {
	c.BeforeHooks = append(c.BeforeHooks, hooks...)
}

//RegisterAfterHook ...
func (c *LoadBalancerHTTPClient) RegisterAfterHook(hooks ...AfterHTTHook) {
	c.AfterHooks = append(c.AfterHooks, hooks...)
}

func (c *LoadBalancerHTTPClient) beforeHTTPHook(ctx context.Context, req *HTTPRequest) {
	for _, h := range c.BeforeHooks {
		h(ctx, req)
	}
}

func (c *LoadBalancerHTTPClient) afterHTTPHook(ctx context.Context, req *HTTPRequest, resp *HTTPResponse, err error) {
	for _, h := range c.AfterHooks {
		h(ctx, req, resp, err)
	}
}
