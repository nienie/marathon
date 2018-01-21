package loadbalancer

import (
	"sync/atomic"

	"github.com/nienie/marathon/server"
)

//RoundRobinRule The most well known and basic load balancing strategy, i.e. Round Robin Rule.
type RoundRobinRule struct {
	BaseRule
	nextServerCyclicCounter int64
}

//NewRoundRobinRule ...
func NewRoundRobinRule() Rule {
	return &RoundRobinRule{}
}

//Choose ...
func (o *RoundRobinRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *RoundRobinRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	if lb == nil {
		return nil
	}

	upList := o.GetLoadBalancer().GetReachableServers()
	upCount := len(upList)
	if upCount == 0 {
		return nil
	}

	nextServerIndex := o.incrementAndGetModulo(upCount)
	return upList[nextServerIndex]
}

func (o *RoundRobinRule) incrementAndGetModulo(modulo int) int {
	for {
		current := atomic.LoadInt64(&o.nextServerCyclicCounter)
		next := (current + 1) % int64(modulo)
		if atomic.CompareAndSwapInt64(&o.nextServerCyclicCounter, current, next) {
			return int(next)
		}
	}
}
