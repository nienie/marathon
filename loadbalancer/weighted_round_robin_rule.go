package loadbalancer

import (
    "sync"
    "sync/atomic"

    "github.com/nienie/marathon/server"
)

//WeightedRoundRobinRule ...
type WeightedRoundRobinRule struct {
    BaseRule
    *sync.RWMutex
    Servers             []*server.Server
    WeightedServerPool  []*server.Server
    Length              int
    isRefreshing            int32
    nextServerCyclicCounter int64
}

//NewWeightedRoundRobinRule ...
func NewWeightedRoundRobinRule() Rule {
    return &WeightedRoundRobinRule{
        RWMutex:        &sync.RWMutex{},
        WeightedServerPool:  make([]*server.Server, 0, 20),
    }
}

//Choose ...
func (o *WeightedRoundRobinRule)Choose(key interface{}) *server.Server {
    return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *WeightedRoundRobinRule)ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
    if lb == nil {
        return nil
    }

    upList := lb.GetReachableServers()
    upCount := len(upList)
    if upCount == 0 {
        return nil
    }

    if !server.CompareServerList(o.Servers, upList) {
        o.RefreshServersAndWeights(upList)
    }

    o.RLock()
    o.RUnlock()
    index := o.incrementAndGetModulo(o.Length)
    return o.WeightedServerPool[index]
}


//SetLoadBalancer ...
func (o *WeightedRoundRobinRule)SetLoadBalancer(lb LoadBalancer) {
    o.BaseRule.SetLoadBalancer(lb)
    if lb != nil {
        o.RefreshServersAndWeights(lb.GetReachableServers())
    }
}

//RefreshServersAndWeights ...
func (o *WeightedRoundRobinRule)RefreshServersAndWeights(servers []*server.Server) {
    if atomic.CompareAndSwapInt32(&o.isRefreshing, int32(0), int32(1)) {
        defer atomic.StoreInt32(&o.isRefreshing, int32(0))
        o.Servers = servers
        o.Lock()
        defer o.Unlock()
        o.WeightedServerPool = o.WeightedServerPool[:0]
        o.Length = 0
        for _, svr := range servers {
            if svr.GetWeight() > 0 {
                for w := 0; w < svr.GetWeight(); w++ {
                    o.WeightedServerPool = append(o.WeightedServerPool, svr)
                    o.Length++
                }
            }
        }
    }
}

func (o *WeightedRoundRobinRule)incrementAndGetModulo(modulo int) int {
    for {
        current := atomic.LoadInt64(&o.nextServerCyclicCounter)
        next := (current + 1) % int64(modulo)
        if atomic.CompareAndSwapInt64(&o.nextServerCyclicCounter, current, next) {
            return int(next)
        }
    }
}