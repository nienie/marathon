package loadbalancer

import (
	"github.com/nienie/marathon/server"
	"math"
)

//LeastResponseTimeRule ...
type LeastResponseTimeRule struct {
	BaseRule
}

//NewLeastResponseTimeRule ...
func NewLeastResponseTimeRule() Rule {
	return &LeastResponseTimeRule{}
}

//Choose ...
func (o *LeastResponseTimeRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *LeastResponseTimeRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	if lb == nil {
		return nil
	}

	allList := o.GetLoadBalancer().GetAllServers()
	totalCount := len(allList)
	if totalCount == 0 {
		return nil
	}

	upList := o.GetLoadBalancer().GetReachableServers()
	upCount := len(upList)
	if upCount == 0 {
		return nil
	}

	lbStats := lb.GetLoadBalancerStats()
	var (
		selectedServer    *server.Server
		leastResponseTime = math.MaxFloat64
	)

	for _, svr := range upList {
		serverStats := lbStats.GetSingleServerStats(svr)
		avgRespTime := serverStats.GetResponseTimeAvg()
		if avgRespTime < leastResponseTime {
			leastResponseTime = avgRespTime
			selectedServer = svr
		}
	}
	return selectedServer
}
