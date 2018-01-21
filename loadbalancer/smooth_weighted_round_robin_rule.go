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

    //upServerList has changed, so refresh the server list and it's weight
    if !server.CompareServerList(o.Servers, upList) {
        o.refreshServersAndWeights(upList)
    }
    o.RLock()
    s := o.Weighted.Next()
    o.RUnlock()
    if s == nil {
        return nil
    }

    return s.(*server.Server)
}

func (o *SmoothWeightedRoundRobinRule)refreshServersAndWeights(servers []*server.Server) {
    if atomic.CompareAndSwapInt32(&o.isRefreshing, int32(0), int32(1)) {
        defer atomic.StoreInt32(&o.isRefreshing, int32(0))
        o.Servers = servers
        o.Lock()
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
        o.refreshServersAndWeights(lb.GetReachableServers())
    }
}