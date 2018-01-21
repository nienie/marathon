package loadbalancer

import (
	"math/rand"
	"time"

	"github.com/nienie/marathon/server"
)

//RandomRule A loadbalacing strategy that randomly distributes traffic amongst existing servers.
type RandomRule struct {
	BaseRule
}

//NewRandomRule ...
func NewRandomRule() Rule {
	return &RandomRule{}
}

//Choose ...
func (o *RandomRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *RandomRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
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

	rand.Seed(time.Now().UnixNano())
	var selectServer *server.Server
	index := rand.Intn(upCount)
	selectServer = upList[index]
	return selectServer
}
