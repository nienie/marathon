package loadbalancer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/loadbalancer/ping"
	"github.com/nienie/marathon/server"
	"github.com/nienie/marathon/utils/timer"
)

const (
	//LoadBalancerPrefix ...
	LoadBalancerPrefix = "LB_"
)

//BaseLoadBalancer ...
type BaseLoadBalancer struct {
	name                  string

	rule                  Rule
	pingAction            ping.Ping
	pingStrategy          ping.Strategy
	lbStats               *Stats

	pingInterval          time.Duration
	recoverInterval       time.Duration

	changeListeners       []server.ListChangeListener
	serverStatusListeners []server.StatusChangeListener

	allServerLock         sync.RWMutex
	upServerLock          sync.RWMutex
	tempDownServerLock    sync.RWMutex
	allServersList        []*server.Server
	upServersList         []*server.Server
	tempDownServerList    []*server.Server

	faultRecoverTimer     *timer.Timer
	healthCheckTimer      *timer.Timer

	pingInProgress        int32
	recoverInProgress     int32
}

//NewBaseLoadBalancer A basic implementation of the load balancer.
func NewBaseLoadBalancer(clientConfig config.ClientConfig, rule Rule, pingAction ping.Ping,
	pingStrategy ping.Strategy) *BaseLoadBalancer {
	loadBalancer := &BaseLoadBalancer{
		name:                  clientConfig.GetClientName(),
		pingAction:            pingAction,
		pingStrategy:          pingStrategy,
		pingInterval:          clientConfig.GetPropertyAsDuration(config.PingInterval, config.DefaultPingInterval),
		recoverInterval:       time.Second * 1,
		changeListeners:       make([]server.ListChangeListener, 0),
		serverStatusListeners: make([]server.StatusChangeListener, 0),
		allServersList:        make([]*server.Server, 0),
		upServersList:         make([]*server.Server, 0),
		tempDownServerList:    make([]*server.Server, 0),
		allServerLock:         sync.RWMutex{},
		upServerLock:          sync.RWMutex{},
		tempDownServerLock:    sync.RWMutex{},
	}
	if loadBalancer.pingStrategy == nil {
		loadBalancer.pingStrategy = ping.NewParallelStrategy()
	}
	loadBalancer.lbStats = NewLoadBalancerStats(clientConfig)
	loadBalancer.SetRule(rule)
	loadBalancer.setupPingTask()
	loadBalancer.setupFaultRecoverTask()
	return loadBalancer
}

//SetName ...
func (o *BaseLoadBalancer) SetName(name string) {
	o.name = name
	o.lbStats.Name = name
}

//GetName ...
func (o *BaseLoadBalancer) GetName() string {
	return o.name
}

//SetRule ...
func (o *BaseLoadBalancer)SetRule(rule Rule) {
	if rule != nil {
		o.rule = rule
	} else {
		o.rule = NewRoundRobinRule()
	}

	o.rule.SetLoadBalancer(o)
}

//GetRule ...
func (o *BaseLoadBalancer)GetRule() Rule {
	return o.rule
}

//SetPing ...
func (o *BaseLoadBalancer)SetPing(ping ping.Ping) {
	o.pingAction = ping
	if ping == nil {
		o.stopPingTask()
		return
	}
	o.setupPingTask()
}

//GetPing ...
func (o *BaseLoadBalancer) GetPing() ping.Ping {
	return o.pingAction
}

func (o *BaseLoadBalancer)setupPingTask() {
	if o.healthCheckTimer != nil {
		o.healthCheckTimer.Cancel()
	}
	o.healthCheckTimer = timer.NewTimer(o.name + "_HealthCheckTask")
	o.healthCheckTimer.Schedule(o.runPingTask, o.pingInterval, 0)
	//o.runPingTask()
}

func (o *BaseLoadBalancer)runPingTask() {
	if !atomic.CompareAndSwapInt32(&o.pingInProgress, 0, 1) {
		return
	}

	defer atomic.StoreInt32(&o.pingInProgress, 0)

	o.allServerLock.RLock()
	allServers := server.CloneServerList(o.allServersList)
	o.allServerLock.RUnlock()

	numCandidates := len(allServers)
	newUpList := make([]*server.Server, 0)
	changeServers := make([]*server.Server, 0)
	pingResults := o.pingStrategy.PingServers(o.pingAction, allServers)

	for i := 0; i < numCandidates; i++ {
		isAlive := pingResults[i]
		svr := allServers[i]
		oldIsAlive := svr.IsAlive()

		svr.SetAlive(isAlive)
		if oldIsAlive != isAlive {
			changeServers = append(changeServers, svr)
		}

		if isAlive {
			newUpList = append(newUpList, svr)
			svr.SetTempDown(false)
		}
	}

	//no servers are alive, make them all be selected.
	if len(newUpList) == 0 {
		o.upServerLock.Lock()
		o.upServersList = allServers
		o.upServerLock.Unlock()
	} else {
		o.upServerLock.Lock()
		o.upServersList = newUpList
		o.upServerLock.Unlock()
	}

	o.notifyServerStatusChangeListener(changeServers)
}

func (o *BaseLoadBalancer) stopPingTask() {
	if o.healthCheckTimer != nil {
		o.healthCheckTimer.Cancel()
	}
}

func (o *BaseLoadBalancer)setupFaultRecoverTask() {
	if o.faultRecoverTimer != nil {
		o.faultRecoverTimer.Cancel()
	}
	o.faultRecoverTimer = timer.NewTimer(o.name + "_FaultRecoverTask")
	o.faultRecoverTimer.Schedule(o.runFaultRecoverTask, o.recoverInterval, 0)
}

func (o *BaseLoadBalancer)runFaultRecoverTask() {
	if !atomic.CompareAndSwapInt32(&o.recoverInProgress, 0, 1) {
		return
	}

	defer atomic.StoreInt32(&o.recoverInProgress, 0)

	newTempDownServers := make([]*server.Server, 0)
	currentTime := time.Duration(time.Now().UnixNano())
	o.tempDownServerLock.Lock()
	for _, svr := range o.tempDownServerList {
		if svr.IsTempDown() == false {
			continue
		}
		if svr.IsAlive() == false {
			svr.SetTempDown(false)
			continue
		}
		stats := o.lbStats.GetSingleServerStats(svr)
		if !stats.IsCircuitBreakerTripped(currentTime) {
			svr.SetTempDown(false)
			continue
		}
		newTempDownServers = append(newTempDownServers, svr)
	}
	o.tempDownServerList = newTempDownServers
	o.tempDownServerLock.Unlock()
}

func (o *BaseLoadBalancer)stopFaultRecoverTask() {
	if o.faultRecoverTimer != nil {
		o.faultRecoverTimer.Cancel()
	}
}

func (o *BaseLoadBalancer) notifyServerStatusChangeListener(changeServes []*server.Server) {
	if changeServes != nil && len(changeServes) != 0 && o.serverStatusListeners != nil {
		for _, serverStatusChangeListener := range o.serverStatusListeners {
			serverStatusChangeListener.ServerStatusChanged(changeServes)
		}
	}
}

//SetPingInterval ...
func (o *BaseLoadBalancer) SetPingInterval(pingInterval time.Duration) {
	if pingInterval < time.Second*1 {
		return
	}

	o.pingInterval = pingInterval
	o.setupPingTask()
}

//GetPingInterval ...
func (o *BaseLoadBalancer) GetPingInterval() time.Duration {
	return o.pingInterval
}

//AddServer Add a server to the 'allServer' list; does not verify uniqueness, so you
// could give a server a greater share by adding it more than once.
func (o *BaseLoadBalancer) AddServer(svr *server.Server) {
	if svr == nil {
		return
	}
	o.allServerLock.RLock()
	numCandidates := len(o.allServersList)
	newList := make([]*server.Server, numCandidates+1)
	copy(newList, o.allServersList)
	o.allServerLock.RUnlock()
	newList[numCandidates] = svr
	o.SetServerList(newList)
}

//SetServerList Set the list of servers used as the server pool. This overrides existing server list.
func (o *BaseLoadBalancer) SetServerList(serverList []*server.Server) {
	var (
		allServers  = make([]*server.Server, 0)
		listChanged bool
	)

	for _, svr := range serverList {
		if svr == nil {
			continue
		}
		allServers = append(allServers, svr)
	}

	if server.CompareServerList(o.allServersList, allServers) {
		listChanged = true
		if o.changeListeners != nil && len(o.changeListeners) > 0 {
			oldList := server.CloneServerList(o.allServersList)
			newList := server.CloneServerList(allServers)
			o.notifyServerListChanged(oldList, newList)
		}
	}
	o.allServerLock.Lock()
	o.allServersList = allServers
	o.allServerLock.Unlock()

	o.tempDownServerLock.Lock()
	o.tempDownServerList = make([]*server.Server, 0)
	o.tempDownServerLock.Unlock()

	if o.pingAction == nil {
		for _, s := range o.allServersList {
			s.SetAlive(true)
		}

		o.upServerLock.Lock()
		o.upServersList = server.CloneServerList(allServers)
		o.upServerLock.Unlock()
		return
	}

	if listChanged {
		o.setupPingTask()
	}
	return
}

func (o *BaseLoadBalancer) notifyServerListChanged(oldList, newList []*server.Server) {
	for _, serverListChangedListener := range o.changeListeners {
		serverListChangedListener.ServerListChanged(oldList, newList)
	}
}

//SetServerListForClusters ...
func (o *BaseLoadBalancer)SetServerListForClusters(clusterServersMap map[string][]*server.Server) {
	o.GetLoadBalancerStats().UpdateClusterServerMapping(clusterServersMap)
}

//AddServers Add a list of servers to the 'allServer' list; does not verify
// uniqueness, so you could give a server a greater share by adding it more than once.
func (o *BaseLoadBalancer) AddServers(servers []*server.Server) {
	if servers == nil || len(servers) == 0 {
		return
	}
	o.allServerLock.RLock()
	allServers := server.CloneServerList(o.allServersList)
	o.allServerLock.RUnlock()
	allServers = append(allServers, servers...)
	o.SetServerList(allServers)
}

//ChooseServer ...
func (o *BaseLoadBalancer) ChooseServer(key interface{}) *server.Server {
	if o.rule == nil {
		return nil
	}
	return o.rule.Choose(key)
}

//MarkServerDown ...
func (o *BaseLoadBalancer) MarkServerDown(svr *server.Server) {
	if svr == nil || svr.IsAlive() == false {
		return
	}
	svr.SetAlive(false)
	o.notifyServerStatusChangeListener([]*server.Server{svr})
}

//GetReachableServers ...
func (o *BaseLoadBalancer) GetReachableServers() []*server.Server {
	o.upServerLock.RLock()
	defer o.upServerLock.RUnlock()
	return server.CloneServerList(o.upServersList)
}

//GetAllServers ...
func (o *BaseLoadBalancer) GetAllServers() []*server.Server {
	o.allServerLock.RLock()
	defer o.allServerLock.RUnlock()
	return server.CloneServerList(o.allServersList)
}

//GetLoadBalancerStats ...
func (o *BaseLoadBalancer) GetLoadBalancerStats() *Stats {
	return o.lbStats
}

//MarkServerTempDown ...
func (o *BaseLoadBalancer)MarkServerTempDown(svr *server.Server) {
	if svr == nil || svr.IsAlive() == false || svr.IsTempDown() == true {
		return
	}
	svr.SetTempDown(true)
	o.tempDownServerLock.Lock()
	o.tempDownServerList = append(o.tempDownServerList, svr)
	o.tempDownServerLock.Unlock()
}

//Shutdown ...
func (o *BaseLoadBalancer) Shutdown() {
	o.stopPingTask()
	o.stopFaultRecoverTask()
}
