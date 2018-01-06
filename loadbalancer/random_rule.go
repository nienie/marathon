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
	var server *server.Server
	for count := 0; count < 20; count++ {
		rand.Seed(time.Now().UnixNano())
		upList := o.GetLoadBalancer().GetReachableServers()
		allList := o.GetLoadBalancer().GetAllServers()
		serverCount := len(allList)
		aliveServerCount := len(upList)

		if serverCount == 0 || aliveServerCount == 0 {
			return nil
		}

		index := rand.Intn(serverCount)
		if index >= aliveServerCount {
			continue
		}
		server = upList[index]
		if server.IsAlive() && server.IsTempDown() == false {
			return server
		}
	}
	return server
}
