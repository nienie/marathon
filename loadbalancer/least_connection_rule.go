package loadbalancer

import (
    "math"

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
func (o *LeastConnectionRule)ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
    if lb == nil {
        return nil
    }

    reachableServers := lb.GetReachableServers()
    allServers := lb.GetAllServers()

    upCount := len(reachableServers)
    serverCount := len(allServers)

    if upCount == 0 || serverCount == 0 {
        return nil
    }

    lbStats := lb.GetLoadBalancerStats()
    var (
        selectedServer *server.Server
        leastConnections int64 = math.MaxInt64
    )
    for _, svr := range reachableServers {
        if svr.IsTempDown() {
            continue
        }
        serverStats := lbStats.GetSingleServerStats(svr)
        cnt := serverStats.GetOpenConnectionsCount()
        if cnt < leastConnections {
            leastConnections = cnt
            selectedServer = svr
        }
    }
    return selectedServer
}