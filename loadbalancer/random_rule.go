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

	upList := o.GetLoadBalancer().GetReachableServers()
	allList := o.GetLoadBalancer().GetAllServers()
	totalCount := len(allList)
	upCount := len(upList)

	if totalCount == 0 || upCount == 0 {
		return nil
	}

	var server *server.Server
	rand.Seed(time.Now().UnixNano())

	for count := 0; count < 20; count++ {
		index := rand.Intn(upCount)

		if index >= upCount {
			continue
		}

		server = upList[index]
		if server.IsAlive() && server.IsTempDown() == false {
			return server
		}
	}

	return server
}
