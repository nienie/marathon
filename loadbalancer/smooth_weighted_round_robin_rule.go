package loadbalancer

import (
    "sync"
    "sync/atomic"

    "github.com/nienie/marathon/server"
    "github.com/smallnest/weighted"
)

//SmoothWeightedRoundRobinRule ...
type SmoothWeightedRoundRobinRule struct {
    BaseRule
    *sync.RWMutex
    isRefreshing int32
    Weighted     *weighted.W1
    Servers      []*server.Server
}

//NewSmoothWeightedRoundRobinRule ...
func NewSmoothWeightedRoundRobinRule() Rule {
    rule := &SmoothWeightedRoundRobinRule{
        Weighted:   &weighted.W1{},
        RWMutex:    &sync.RWMutex{},
    }
    return rule
}

//Choose ...
func (o *SmoothWeightedRoundRobinRule)Choose(key interface{}) *server.Server {
    return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *SmoothWeightedRoundRobinRule)ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
    if lb == nil {
        return nil
    }

    upList := lb.GetReachableServers()
    upCount := len(upList)
    if upCount == 0 {
        return nil
    }

    o.RLock()
    isEqual := server.CompareServerList(o.Servers, upList)
    o.RUnlock()
    //upServerList has changed, so refresh the server list and weights
    if !isEqual {
        o.RefreshServersAndWeights(upList)
    }
    o.Lock()
    s := o.Weighted.Next()
    o.Unlock()
    if s == nil {
        return nil
    }

    return s.(*server.Server)
}

//RefreshServersAndWeights ...
func (o *SmoothWeightedRoundRobinRule)RefreshServersAndWeights(servers []*server.Server) {
    if atomic.CompareAndSwapInt32(&o.isRefreshing, int32(0), int32(1)) {
        defer atomic.StoreInt32(&o.isRefreshing, int32(0))
        o.Lock()
        o.Servers = servers
        o.Weighted.RemoveAll()
        for _, svr := range o.Servers {
            o.Weighted.Add(svr, svr.Weight)
        }
        o.Unlock()
    }
}

//SetLoadBalancer ...
func (o *SmoothWeightedRoundRobinRule)SetLoadBalancer(lb LoadBalancer) {
    o.BaseRule.SetLoadBalancer(lb)
    if lb != nil {
        o.RefreshServersAndWeights(lb.GetReachableServers())
    }
}