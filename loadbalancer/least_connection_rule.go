package loadbalancer

import (
	"math"
	"time"

	"github.com/nienie/marathon/server"
)

//LeastConnectionRule ...
type LeastConnectionRule struct {
	BaseRule
}

//NewLeastConnectionRule ...
func NewLeastConnectionRule() Rule {
	return &LeastConnectionRule{}
}

//Choose ...
func (o *LeastConnectionRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *LeastConnectionRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	if lb == nil {
		return nil
	}

	upList := o.GetLoadBalancer().GetReachableServers()
	upCount := len(upList)
	if upCount == 0 {
		return nil
	}

	lbStats := lb.GetLoadBalancerStats()
	var (
		selectedServer   *server.Server
		leastConnections int64 = math.MaxInt64
	)
	currentTime := time.Duration(time.Now().UnixNano())
	for _, svr := range upList {
		serverStats := lbStats.GetSingleServerStats(svr)
		cnt := serverStats.GetActiveRequestsCount(currentTime)
		if cnt < leastConnections {
			leastConnections = cnt
			selectedServer = svr
		}
	}
	return selectedServer
}
