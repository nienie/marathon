package client

import (
	"net/url"
)

//Request An object that represents a common client request that is suitable for all communication protocol.
type Request interface {
	//GetURI ...
	GetURI() *url.URL

	//GetLoadBalancerKey ...
	GetLoadBalancerKey() interface{}

	//GetIsRetriable ...
	//GetIsRetriable() bool

	ReplaceURI(*url.URL)
}
