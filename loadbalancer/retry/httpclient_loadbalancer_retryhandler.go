package retry

import (
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/errors"
)

//HTTPClientLoadBalancerRetryHandler ...
type HTTPClientLoadBalancerRetryHandler struct {
	*LoadBalancerRetryHandler
}

//NewDefaultHTTPClientLoadBalancerRetryHandler ...
func NewDefaultHTTPClientLoadBalancerRetryHandler() Handler {
	return &HTTPClientLoadBalancerRetryHandler{
		NewDefaultLoadBalancerRetryHandler(),
	}
}

//NewHTTPClientLoadBalancerRetryHandler ...
func NewHTTPClientLoadBalancerRetryHandler(clientConfig config.ClientConfig) Handler {
	return &HTTPClientLoadBalancerRetryHandler{
		NewLoadBalancerRetryHandler(clientConfig),
	}
}

//IsCircuitTrippingException ...
func (c *HTTPClientLoadBalancerRetryHandler) IsCircuitTrippingException(err error) bool {
	switch err.(type) {
	case errors.ClientError:
		errorType := err.(errors.ClientError).GetErrType()
		switch errorType {
		case errors.ServerThrottled:
			return true
		case errors.UnknownHostException:
			return true
		default:
			return c.LoadBalancerRetryHandler.IsCircuitTrippingException(err)
		}
	default:
		return c.LoadBalancerRetryHandler.IsCircuitTrippingException(err)
	}
}

//IsRetriableException ...
func (c *HTTPClientLoadBalancerRetryHandler) IsRetriableException(err error, sameServer bool) bool {
	switch err.(type) {
	case errors.ClientError:
		errorType := err.(errors.ClientError).GetErrType()
		if errorType == errors.ServerThrottled {
			return !sameServer && c.RetryEnabled
		}
	default:
	}
	return c.LoadBalancerRetryHandler.IsCircuitTrippingException(err)
}
