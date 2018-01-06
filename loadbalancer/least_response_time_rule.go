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
func (o *LeastResponseTimeRule)Choose(key interface{}) *server.Server {
    return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *LeastResponseTimeRule)ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
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
        leastResponseTime = math.MaxFloat64
    )

    for _, svr := range reachableServers {
        if svr.IsTempDown() {
            continue
        }
        serverStats := lbStats.GetSingleServerStats(svr)
        avg := serverStats.GetResponseTimeAvg()
        if avg < leastResponseTime {
            leastResponseTime = avg
            selectedServer = svr
        }
    }
    return selectedServer
}

