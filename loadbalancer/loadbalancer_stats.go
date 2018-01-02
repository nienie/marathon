package loadbalancer

import (
	"fmt"
	"sync"
	"time"

	"github.com/nienie/marathon/server"
	"github.com/nienie/marathon/utils/cache"
	"github.com/nienie/marathon/config"
)

const (
	//DefaultConnectionFailureThreshold ...
	DefaultConnectionFailureThreshold = 3
	//DefaultCircuitTrippedTimeoutFactor ...
	DefaultCircuitTrippedTimeoutFactor = 10
	//DefaultCircuitTripMaxTimeout ...
	DefaultCircuitTripMaxTimeout = 30 * time.Second
	//DefaultServerStatsCacheExpireTime ...
	DefaultServerStatsCacheExpireTime = 30 * time.Minute
)

var (
	errWrongType = fmt.Errorf("wrong type")
)

//Stats Class that acts as a repository of operational charateristics and statistics
//of every Node/Server in the LaodBalancer. This information can be used to just observe and understand the runtime
//behavior of the loadbalancer or more importantly for the basis that determines the loadbalacing strategy
type Stats struct {
	Name                          string
	ConnectionFailureThreshold    int
	CircuitTrippedTimeoutFactor   int
	MaxCircuitTrippedTimeout      time.Duration
	ServerStatsCacheExpireTime    time.Duration

	serverStatsCache              *cache.TimedCache
	clusterStatsMap			      map[string]*ClusterStats
	clusterStatsLock 			  sync.RWMutex
	upServerClusterMap			  map[string][]*server.Server
	serverClusterLock 			  sync.RWMutex
}

type serverStatsCacheCallback struct {
	loadBalancerStats *Stats
}

func newServerStatsCacheCallback(stats *Stats) cache.Callback {
	return &serverStatsCacheCallback{
		loadBalancerStats: stats,
	}
}

//OnLoad ...
func (o *serverStatsCacheCallback) OnLoad(key interface{}) (interface{}, error) {
	server, ok := key.(*server.Server)
	if !ok {
		return nil, errWrongType
	}
	return o.loadBalancerStats.CreateServerStats(server), nil
}

//OnRemove ...
func (o *serverStatsCacheCallback) OnRemove(key interface{}, val interface{}) error {
	serverStats, ok := val.(*server.Stats)
	if !ok {
		return errWrongType
	}
	serverStats.Close()
	return nil
}

//NewLoadBalancerStats ...
func NewLoadBalancerStats(clientConfig config.ClientConfig) *Stats {
	loadBalancerStats := &Stats{
		Name: 						   clientConfig.GetClientName(),
		ConnectionFailureThreshold:    clientConfig.GetPropertyAsInteger(config.ConnectionFailureThreshold,
			config.DefaultConnectionFailureThreshold),
		CircuitTrippedTimeoutFactor:   clientConfig.GetPropertyAsInteger(config.CircuitTrippedTimeoutFactor,
			config.DefaultCircuitTrippedTimeoutFactor),
		MaxCircuitTrippedTimeout:      clientConfig.GetPropertyAsDuration(config.CircuitTripMaxTimeout,
			config.DefaultCircuitTripMaxTimeout),
		ServerStatsCacheExpireTime:    DefaultServerStatsCacheExpireTime,
		clusterStatsMap:			   make(map[string]*ClusterStats),
		clusterStatsLock:			   sync.RWMutex{},
		upServerClusterMap:			   make(map[string][]*server.Server),
		serverClusterLock:			   sync.RWMutex{},
	}
	callback := newServerStatsCacheCallback(loadBalancerStats)
	loadBalancerStats.serverStatsCache = cache.NewTimedCache(loadBalancerStats.ServerStatsCacheExpireTime, callback)
	return loadBalancerStats
}

//CreateServerStats ...
func (o *Stats) CreateServerStats(s *server.Server) *server.Stats {
	ss := server.NewDefaultServerStats()
	ss.CircuitTrippedTimeoutFactor = o.CircuitTrippedTimeoutFactor
	ss.ConnectionFailureThreshold = o.ConnectionFailureThreshold
	ss.MaxCircuitTrippedTimeout = o.MaxCircuitTrippedTimeout
	ss.Initialize(s)
	return ss
}

//UpdateServerList The caller o this class is tasked to call this method every so often if
//the servers participating in the LoadBalancer changes
func (o *Stats) UpdateServerList(servers []*server.Server) {
	for _, server := range servers {
		o.AddServer(server)
	}
}

//AddServer ...
func (o *Stats) AddServer(svr *server.Server) *server.Stats {
	return o.GetSingleServerStats(svr)
}

//GetSingleServerStats ...
func (o *Stats) GetSingleServerStats(svr *server.Server) *server.Stats {
	if svr == nil {
		return nil
	}
	ss, err := o.serverStatsCache.GetAndSetWhenNotExisted(svr)
	if err != nil || ss == nil {
		serverStats := o.CreateServerStats(svr)
		o.serverStatsCache.Set(svr, serverStats, time.Duration(0))
		return serverStats
	}
	serverStats, ok := ss.(*server.Stats)
	if !ok {
		serverStats = o.CreateServerStats(svr)
		o.serverStatsCache.Set(svr, serverStats, time.Duration(0))
		return serverStats
	}
	return serverStats
}

//NoteResponseTime ...
func (o *Stats) NoteResponseTime(server *server.Server, msec float64) {
	ss := o.GetSingleServerStats(server)
	ss.NoteResponseTime(msec)
}

//IncrementActiveRequestsCount ...
func (o *Stats) IncrementActiveRequestsCount(server *server.Server) {
	ss := o.GetSingleServerStats(server)
	ss.IncrementActiveRequestsCount()
}

//DecrementActiveRequestsCount ...
func (o *Stats) DecrementActiveRequestsCount(server *server.Server) {
	ss := o.GetSingleServerStats(server)
	ss.DecrementActiveRequestsCount()
}

//IsCircuitBreakerTripped ...
func (o *Stats) IsCircuitBreakerTripped(server *server.Server) bool {
	ss := o.GetSingleServerStats(server)
	return ss.IsCircuitBreakerTripped(time.Duration(time.Now().UnixNano()))
}

//IncrementSuccessiveConnectionFailureCount ...
func (o *Stats) IncrementSuccessiveConnectionFailureCount(server *server.Server) {
	ss := o.GetSingleServerStats(server)
	ss.IncrementSuccessiveConnectionFailureCount()
}

//ClearSuccessiveConnectionFailureCount ...
func (o *Stats) ClearSuccessiveConnectionFailureCount(server *server.Server) {
	ss := o.GetSingleServerStats(server)
	ss.ClearSuccessiveConnectionFailureCount()
}

//IncrementNumRequests ...
func (o *Stats) IncrementNumRequests(server *server.Server) {
	ss := o.GetSingleServerStats(server)
	ss.IncrementNumRequests()
}

//GetAllServerStats ...
func (o *Stats) GetAllServerStats() map[*server.Server]*server.Stats {
	m := o.serverStatsCache.ToMap()
	serverStatsMap := make(map[*server.Server]*server.Stats)
	for s, ss := range m {
		svr, ok := s.(*server.Server)
		if !ok {
			continue
		}
		svrStats, ok := ss.(*server.Stats)
		if !ok {
			continue
		}
		serverStatsMap[svr] = svrStats
	}
	return serverStatsMap
}

//GetClusterSnapshotByName ...
func (o *Stats)GetClusterSnapshotByName(cluster string) *ClusterSnapshot {
	if len(cluster) == 0 {
		return NewDefaultClusterSnapshot()
	}
	o.serverClusterLock.RLock()
	currentList := o.upServerClusterMap[cluster]
	o.serverClusterLock.RUnlock()
	return o.GetClusterSnapshotByServers(currentList)
}

//GetClusterSnapshotByServers ...
func (o *Stats)GetClusterSnapshotByServers(servers []*server.Server) *ClusterSnapshot {
	if servers == nil || len(servers) == 0 {
		return NewDefaultClusterSnapshot()
	}
	var (
		instanceCount = len(servers)
		activeConnectionsCount int64
		activeConnectionsCountOnAvailableServer int64
		circuitBreakerTrippedCount int
		loadPerServer float64
		currentTime = time.Duration(time.Now().UnixNano())
	)
	for _, svr := range servers {
		stat := o.GetSingleServerStats(svr)
		if stat.IsCircuitBreakerTripped(currentTime) {
			circuitBreakerTrippedCount++
		} else {
			activeConnectionsCountOnAvailableServer += stat.GetActiveRequestsCount(currentTime)
		}
		activeConnectionsCount += stat.GetActiveRequestsCount(currentTime)
	}
	if circuitBreakerTrippedCount == instanceCount {
		if instanceCount > 0 {
			loadPerServer = -1
		}
	} else {
		loadPerServer = float64(activeConnectionsCountOnAvailableServer) / float64(instanceCount - circuitBreakerTrippedCount)
	}
	return NewClusterSnapshot(instanceCount, loadPerServer, circuitBreakerTrippedCount, activeConnectionsCount)
}

//GetInstanceCount ...
func (o *Stats)GetInstanceCount(cluster string) int {
	if len(cluster) == 0 {
		return 0
	}
	o.serverClusterLock.RLock()
	currentList := o.upServerClusterMap[cluster]
	o.serverClusterLock.RUnlock()
	if currentList == nil {
		return 0
	}
	return len(currentList)
}

//GetActiveRequestsCount ...
func (o *Stats)GetActiveRequestsCount(cluster string) int64 {
	snapshot := o.GetClusterSnapshotByName(cluster)
	return snapshot.ActiveRequestsCount
}

//GetActiveRequestsPerServer ...
func (o *Stats)GetActiveRequestsPerServer(cluster string) float64 {
	snapshot := o.GetClusterSnapshotByName(cluster)
	return snapshot.LoadPerServer
}

//GetCircuitBreakerTrippedCount ...
func (o *Stats)GetCircuitBreakerTrippedCount(cluster string) int {
	snapshot := o.GetClusterSnapshotByName(cluster)
	return snapshot.CircuitTrippedCount
}

//GetMeasuredClusterHits ...
func (o *Stats)GetMeasuredClusterHits(cluster string) int64 {
	if len(cluster) == 0 {
		return 0
	}
	var count int64
	o.serverClusterLock.RLock()
	currentList, ok := o.upServerClusterMap[cluster]
	o.serverClusterLock.RUnlock()
	if !ok {
		return 0
	}
	for _, svr := range currentList {
		stat := o.GetSingleServerStats(svr)
		count += stat.GetMeasuredRequestsCount()
	}
	return count
}

//GetAvailableClusters ...
func (o *Stats)GetAvailableClusters() []string {
	clusters := make([]string, 0)
	o.serverClusterLock.RLock()
	for cluster := range o.upServerClusterMap {
		clusters = append(clusters, cluster)
	}
	o.serverClusterLock.RUnlock()
	return clusters
}

func (o *Stats)getClusterStats(cluster string) *ClusterStats {
	o.clusterStatsLock.Lock()
	clusterStats := o.clusterStatsMap[cluster]
	if clusterStats == nil {
		clusterStats := NewClusterStats(cluster, o)
		o.clusterStatsMap[cluster] = clusterStats
	}
	o.clusterStatsLock.Unlock()
	return clusterStats
}

//UpdateClusterServerMapping ...
func (o *Stats)UpdateClusterServerMapping(m map[string][]*server.Server) {
	newMap := make(map[string][]*server.Server)
	clusters := make([]string, 0)
	for key, val := range m {
		clusters = append(clusters, key)
		newMap[key] = val
	}
	o.clusterStatsLock.Lock()
	o.upServerClusterMap = newMap
	o.clusterStatsLock.Unlock()
	for _, cluster := range clusters {
		o.getClusterStats(cluster)
	}
}
