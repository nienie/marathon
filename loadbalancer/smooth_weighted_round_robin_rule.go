package loadbalancer

import (
    "github.com/nienie/marathon/server"
    "github.com/smallnest/weighted"
)

//SmoothWeightedRoundRobinRule ...
type SmoothWeightedRoundRobinRule struct {
    BaseRule
    Weighted *weighted.W1
    Servers []*server.Server
}

//NewSmoothWeightedRoundRobinRule ...
func NewSmoothWeightedRoundRobinRule() Rule {
    rule := &SmoothWeightedRoundRobinRule{
        Weighted:  &weighted.W1{},
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

    s := o.Weighted.Next()
    if s == nil {
        return nil
    }

    return s.(*server.Server)
}

func (o *SmoothWeightedRoundRobinRule)refreshServersAndWeights(servers []*server.Server) {
    o.Servers = servers
    o.Weighted.RemoveAll()
    for _, svr := range o.Servers {
        o.Weighted.Add(svr, svr.Weight)
    }
}

//SetLoadBalancer ...
func (o *SmoothWeightedRoundRobinRule)SetLoadBalancer(lb LoadBalancer) {
    o.BaseRule.SetLoadBalancer(lb)
    if lb != nil {
        o.refreshServersAndWeights(lb.GetReachableServers())
    }
}