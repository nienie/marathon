package loadbalancer

import (
	"github.com/nienie/marathon/server"
)

//Rule Interface that defines a "Rule" for a LoadBalancer. A Rule can be thought of
//as a Strategy for loadbalacing. Well known loadbalancing strategies include
//Round Robin, Response Time based etc.
type Rule interface {
	//Choose choose one alive server from lb.allServers or lb.upServers according to key
	Choose(key interface{}) *server.Server

	//SetLoadBalancer ...
	SetLoadBalancer(lb LoadBalancer)

	//GetLoadBalancer ...
	GetLoadBalancer() LoadBalancer
}

//BaseRule class that provides a default implementation for setting and getting load balancer
type BaseRule struct {
	lb LoadBalancer
}

//SetLoadBalancer ...
func (o *BaseRule) SetLoadBalancer(lb LoadBalancer) {
	o.lb = lb
}

//GetLoadBalancer ...
func (o *BaseRule) GetLoadBalancer() LoadBalancer {
	return o.lb
}

//Choose ...
func (o *BaseRule) Choose(key interface{}) *server.Server {
	return nil
}
