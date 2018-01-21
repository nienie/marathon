package marathon

import (
	"sync"

	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/httpclient"
	"github.com/nienie/marathon/loadbalancer"
	"github.com/nienie/marathon/logger"
	"github.com/nienie/marathon/loadbalancer/ping"
)

var (
	cf *clientFactory
	ruleMap map[string]RuleConstructor
	pingStrategyMap map[string]PingStrategyConstructor
)

//RuleConstructor ...
type RuleConstructor func() loadbalancer.Rule

//PingStrategyConstructor ...
type PingStrategyConstructor func() ping.Strategy

func init() {
	ruleMap = map[string]RuleConstructor{
		config.SmoothWeightedRoundRobinRule: 	func() loadbalancer.Rule {
			return loadbalancer.NewSmoothWeightedRoundRobinRule()
		},
		config.WeightedRoundRobinRule:		func() loadbalancer.Rule {
			return loadbalancer.NewWeightedRoundRobinRule()
		},
		config.RoundRobinRule:		func() loadbalancer.Rule {
			return loadbalancer.NewRoundRobinRule()
		},
		config.HashRule:		func() loadbalancer.Rule {
			return loadbalancer.NewHashRule()
		},
		config.RandomRule:		func() loadbalancer.Rule {
			return loadbalancer.NewRandomRule()
		},
		config.LeastConnectionRule:		func() loadbalancer.Rule{
			return loadbalancer.NewLeastConnectionRule()
		},
		config.LeastResponseTimeRule:		func() loadbalancer.Rule {
			return loadbalancer.NewLeastResponseTimeRule()
		},
		config.WeightedResponseTimeRule:		func() loadbalancer.Rule {
			return loadbalancer.NewWeightedResponseTimeRule()
		},
	}
	pingStrategyMap = map[string]PingStrategyConstructor {
		config.ParallelPingStrategy:		func()ping.Strategy {
			return ping.NewParallelStrategy()
		},
		config.SerialPingStrategy:		func()ping.Strategy {
			return ping.NewSerialStrategy()
		},
	}
	cf = newClientFactory()
}

type clientFactory struct {
	loadBalancers map[string]loadbalancer.LoadBalancer
	lbLock        *sync.RWMutex
	clients       map[string]*httpclient.LoadBalancerHTTPClient
	clientLock    *sync.RWMutex
}

func newClientFactory() *clientFactory {
	return &clientFactory{
		loadBalancers: make(map[string]loadbalancer.LoadBalancer),
		lbLock:        &sync.RWMutex{},
		clients:       make(map[string]*httpclient.LoadBalancerHTTPClient),
		clientLock:    &sync.RWMutex{},
	}
}

//RegisterLoadBalancer ...
func RegisterLoadBalancer(name string, lb loadbalancer.LoadBalancer) {
	if len(name) == 0 || lb == nil {
		return
	}
	if _, ok := cf.loadBalancers[name]; ok {
		return
	}
	cf.lbLock.Lock()
	if _, ok := cf.loadBalancers[name]; !ok {
		cf.loadBalancers[name] = lb
	}
	cf.lbLock.Unlock()
}

//RegisterHTTPClient ...
func RegisterHTTPClient(name string, client *httpclient.LoadBalancerHTTPClient) {
	if len(name) == 0 || client == nil {
		return
	}
	if _, ok := cf.clients[name]; ok {
		return
	}
	cf.clientLock.Lock()
	if _, ok := cf.clients[name]; !ok {
		cf.clients[name] = client
	}
	cf.clientLock.Unlock()
}

//GetLoadBalancerByName ...
func GetLoadBalancerByName(name string) loadbalancer.LoadBalancer {
	return cf.loadBalancers[name]
}

//GetHTTPClientByName ...
func GetHTTPClientByName(name string) *httpclient.LoadBalancerHTTPClient {
	return cf.clients[name]
}

//SetLogger ...
func SetLogger(l logger.Logger) {
	logger.SetLogger(l)
}

//GetBaseLoadBalancer ...
func GetBaseLoadBalancer(clientConfig config.ClientConfig) loadbalancer.LoadBalancer {
	clientName := clientConfig.GetClientName()
	lb := GetLoadBalancerByName(clientName)
	if lb != nil {
		return lb
	}

	ruleName := clientConfig.GetPropertyAsString(config.LoadBalancerRule, config.SmoothWeightedRoundRobinRule)
	if _, ok := ruleMap[ruleName]; !ok {
		ruleName = config.SmoothWeightedRoundRobinRule
	}
	rule := ruleMap[ruleName]()

	pingStrategyName := clientConfig.GetPropertyAsString(config.PingStrategy, config.ParallelPingStrategy)
	if _, ok := pingStrategyMap[pingStrategyName]; !ok {
		pingStrategyName = config.ParallelPingStrategy
	}
	strategy := pingStrategyMap[pingStrategyName]()

	lb = loadbalancer.NewBaseLoadBalancer(clientConfig, rule, nil, strategy)
	cf.lbLock.Lock()
	cf.loadBalancers[clientName] = lb
	cf.lbLock.Unlock()
	return lb
}

//GetHTTPClient ...
func GetHTTPClient(clientConfig config.ClientConfig) *httpclient.LoadBalancerHTTPClient {
	clientName := clientConfig.GetClientName()
	client := GetHTTPClientByName(clientName)
	if client != nil {
		return client
	}
	client = httpclient.NewHTTPLoadBalancerClient(clientConfig, nil)
	cf.clientLock.Lock()
	cf.clients[clientName] = client
	cf.clientLock.Unlock()
	return client
}
