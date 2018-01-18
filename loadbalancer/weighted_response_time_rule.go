package loadbalancer

import (
    "math"
    "sync"

    "github.com/nienie/marathon/server"
)

//WeightedResponseTimeRule ...
type WeightedResponseTimeRule struct {
    BaseRule
}

//NewWeightedResponseTimeRule ...
func NewWeightedResponseTimeRule() Rule {
    return &WeightedResponseTimeRule{}
}

//Choose ...
func (r *WeightedResponseTimeRule)Choose(key interface{}) *server.Server {
    return r.ChooseFromLoadBalancer(r.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (r *WeightedResponseTimeRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
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

    var (
        wg = &sync.WaitGroup{}
        serversResponseTime = make([]float64, upCount)
        leastResponseTime = math.MaxFloat64
        selectedServer *server.Server
    )

    for i, svr := range reachableServers {
        wg.Add(1)
        serverStats := lb.GetLoadBalancerStats().GetSingleServerStats(svr)
        go func (serverStats *server.Stats, index int) {
            defer wg.Done()
            avgRespTimePerSeconds := serverStats.GetAvgResponseTimePerSecond()
            l := len(avgRespTimePerSeconds)
            if l == 0 {
                serversResponseTime[index] = 0
                return
            }

            if l == 1 {
                serversResponseTime[index] = avgRespTimePerSeconds[0]
                return
            }
            delta := float64(1 / (l - 1))
            for j, avgRespTimePerSecond := range avgRespTimePerSeconds {
                serversResponseTime[index] = (0.5 + float64(j) * delta) * avgRespTimePerSecond
            }
        }(serverStats, i)
    }
    wg.Wait()

    for i, svr := range reachableServers {
        if serversResponseTime[i] < leastResponseTime && svr.IsTempDown() == false {
            leastResponseTime = serversResponseTime[i]
            selectedServer = svr
        }
    }

    return selectedServer
}