package loadbalancer

import (
	"github.com/nienie/marathon/server"
)

//LoadBalancer Interface that defines the operations for a software loadbalancer.
// A typical loadbalancer minimally need a set of servers to loadbalance for,
// a method to mark a particular server to be out of rotation and a call that
// will choose a server from the existing list of server.
type LoadBalancer interface {
	//AddServers initializes list of servers.
	//This API also serves to add additional ones at a later time.
	//The same logical server (host:port) could essentially be added multiple times
	//(helpful in cases where you want to give more "weightage" perhaps ..)
	AddServers(newServers []*server.Server)

	//ChooseServer choosees a server from load balancer.
	//key is an interface that the load balancer may use to determin which server to return. nill if
	// the load balancer does not use this parameter.
	ChooseServer(key interface{}) *server.Server

	//MarkServerDown to be called by the clients of the load balancer to notify that a Server is down
	//else, the load balancer will think its still Alive until the next Ping cycle.
	MarkServerDown(server *server.Server)

	//MarkServerTempDown mark a server down temporary ...
	MarkServerTempDown(*server.Server)

	//GetReachableServers only the servers that are up and reachable.
	GetReachableServers() []*server.Server

	//GetAllServers all known servers, both reachable and unreachable.
	GetAllServers() []*server.Server

	//GetLoadBalancerStats ...
	GetLoadBalancerStats() *Stats
}
