package loadbalancer

import (
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/server"
	"sync/atomic"
)

//DynamicServerListLoadBalancer ...
type DynamicServerListLoadBalancer struct {
	*BaseLoadBalancer
	Filter            server.ListFilter
	ServerListImp     server.List
	ServerListUpdater server.ListUpdater

	serverListUpdaterInProgress int32
}

//NewDynamicServerListLoadBalancer ...
func NewDynamicServerListLoadBalancer(clientConfig config.ClientConfig, rule Rule, serverListImp server.List) *DynamicServerListLoadBalancer {
	lb := &DynamicServerListLoadBalancer{
		BaseLoadBalancer: NewBaseLoadBalancer(clientConfig, rule, nil, nil),
		ServerListImp:    serverListImp,
		ServerListUpdater: server.NewPollingServerListUpdater(
			clientConfig.GetPropertyAsDuration(config.ListOfServersPollingInterval, config.DefaultListOfServersPollingInterval),
		),
		Filter: nil,
		serverListUpdaterInProgress: int32(0),
	}
	lb.init()
	return lb
}

func (o *DynamicServerListLoadBalancer) init() {
	o.ServerListUpdater.Start(o)
	o.UpdateListOfServers()
}

//UpdateListOfServers ...
func (o *DynamicServerListLoadBalancer) UpdateListOfServers() {
	if o.ServerListImp != nil {
		servers := o.ServerListImp.GetUpdatedListOfServers()

		if o.Filter != nil {
			servers = o.Filter.GetFilteredListOfServers(servers)
		}
		o.UpdateAllServerList(servers)
	}
}

//UpdateAllServerList ...
func (o *DynamicServerListLoadBalancer) UpdateAllServerList(servers []*server.Server) {
	if atomic.CompareAndSwapInt32(&o.serverListUpdaterInProgress, 0, 1) {
		defer atomic.StoreInt32(&o.serverListUpdaterInProgress, 0)
		for _, svr := range servers {
			svr.SetAlive(true)
			svr.SetTempDown(false)
		}
		o.SetServerList(servers)
		o.runPingTask()
	}
}

//SetServerList ...
func (o *DynamicServerListLoadBalancer) SetServerList(servers []*server.Server) {
	o.BaseLoadBalancer.SetServerList(servers)
	serversInClusters := make(map[string][]*server.Server)
	for _, svr := range servers {
		o.GetLoadBalancerStats().GetSingleServerStats(svr)
		cluster := svr.GetCluster()
		if len(cluster) > 0 {
			clusterServers, ok := serversInClusters[cluster]
			if !ok {
				clusterServers = make([]*server.Server, 0)
				serversInClusters[cluster] = clusterServers
			}
			clusterServers = append(clusterServers, svr)
		}
	}
	o.SetServerListForClusters(serversInClusters)
}

//DoUpdate ...
func (o *DynamicServerListLoadBalancer) DoUpdate() {
	o.UpdateListOfServers()
}

//StopServerListRefreshing ...
func (o *DynamicServerListLoadBalancer) StopServerListRefreshing() {
	if o.ServerListUpdater != nil {
		o.ServerListUpdater.Stop()
	}
}

//Shutdown ...
func (o *DynamicServerListLoadBalancer) Shutdown() {
	o.BaseLoadBalancer.Shutdown()
	o.StopServerListRefreshing()
}
