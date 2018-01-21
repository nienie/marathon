package loadbalancer

import (
	"math/rand"
	"time"

	"github.com/nienie/marathon/server"
)

//RandomRule A loadbalacing strategy that randomly distributes traffic amongst existing servers.
type RandomRule struct {
	BaseRule
	Random *rand.Rand
}

//NewRandomRule ...
func NewRandomRule() Rule {
	s := rand.NewSource(time.Now().UnixNano())
	return &RandomRule{
		Random:		rand.New(s),
	}
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

	var selectServer *server.Server
	index := o.Random.Intn(upCount)
	selectServer = upList[index]
	return selectServer
}
