package retry

import (
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/errors"
)

//LoadBalancerRetryHandler a default implementation.
type LoadBalancerRetryHandler struct {
	RetrySameServer int
	RetryNextServer int
	RetryEnabled    bool
}

//NewLoadBalancerRetryHandler ...
func NewLoadBalancerRetryHandler(clientConfig config.ClientConfig) *LoadBalancerRetryHandler {
	if clientConfig == nil {
		return NewDefaultLoadBalancerRetryHandler()
	}
	return &LoadBalancerRetryHandler{
		RetryEnabled:    clientConfig.GetPropertyAsBool(config.OKToRetryOnAllOperations, config.DefaultOKToRetryOnAllOperations),
		RetrySameServer: clientConfig.GetPropertyAsInteger(config.MaxAutoRetries, config.DefaultMaxAutoRetries),
		RetryNextServer: clientConfig.GetPropertyAsInteger(config.MaxAutoRetriesNextServer, config.DefaultMaxAutoRetriesNextServer),
	}
}

//NewDefaultLoadBalancerRetryHandler ...
func NewDefaultLoadBalancerRetryHandler() *LoadBalancerRetryHandler {
	return &LoadBalancerRetryHandler{
		RetryEnabled:    config.DefaultOKToRetryOnAllOperations,
		RetrySameServer: config.DefaultMaxAutoRetries,
		RetryNextServer: config.DefaultMaxAutoRetriesNextServer,
	}
}

//IsRetriableException ...
func (o *LoadBalancerRetryHandler) IsRetriableException(err error, sameServer bool) bool {
	if o.RetryEnabled {
		if sameServer {
			switch err.(type) {
			case errors.ClientError:
				errorType := err.(errors.ClientError).GetErrType()
				if errorType == errors.SocketTimeoutException || errorType == errors.ConnectException {
					return true
				}
				return false
			default:
				return false
			}
		} else {
			return true
		}
	}
	return false
}

//IsCircuitTrippingException ...
func (o *LoadBalancerRetryHandler) IsCircuitTrippingException(err error) bool {
	switch err.(type) {
	case errors.ClientError:
		errorType := err.(errors.ClientError).GetErrType()
		if errorType == errors.SocketTimeoutException || errorType == errors.SocketException {
			return true
		}
		return false
	default:
		return false
	}
}

//GetMaxRetriesOnSameServer ...
func (o *LoadBalancerRetryHandler) GetMaxRetriesOnSameServer() int {
	return o.RetrySameServer
}

//GetMaxRetriesOnNextServer ...
func (o *LoadBalancerRetryHandler) GetMaxRetriesOnNextServer() int {
	return o.RetryNextServer
}
