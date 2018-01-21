package ping

import (
	"sync"

	"github.com/nienie/marathon/logger"
	"github.com/nienie/marathon/server"
)

//Strategy defines the strategy, used to ping all servers, registered in BaseLoadBalancer.
//You would typically create custom implementation of this interface, if you
//want your servers to be pinged in parallel.
type Strategy interface {
	//PingServers ...
	PingServers(ping Ping, servers []*server.Server) []bool
}

//SerialStrategy performs ping serially
type SerialStrategy struct{}

//NewSerialStrategy ...
func NewSerialStrategy() Strategy {
	return &SerialStrategy{}
}

//PingServers ...
func (o *SerialStrategy) PingServers(ping Ping, servers []*server.Server) []bool {
	numCandidates := len(servers)
	results := make([]bool, numCandidates)

	for i := 0; i < numCandidates; i++ {
		results[i] = tryPing(ping, servers[i])
	}
	return results
}

//ParallelStrategy performs ping concurrently.
type ParallelStrategy struct{}

//NewParallelStrategy ...
func NewParallelStrategy() Strategy {
	return &ParallelStrategy{}
}

//PingServers ...
func (o *ParallelStrategy) PingServers(ping Ping, servers []*server.Server) []bool {
	numCandidates := len(servers)
	results := make([]bool, numCandidates)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCandidates; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = tryPing(ping, servers[i])
		}()
	}
	wg.Wait()
	return results
}

func tryPing(ping Ping, server *server.Server) bool {
	isAlive := false
	defer func() {
		if r := recover(); r != nil {
			logger.Warnf(nil, "err_msg=ping server %s paniced||err=%v", server.GetHostPort(), r)
			return
		}
	}()
	if ping != nil {
		isAlive = ping.IsAlive(server)
	}
	return isAlive
}
