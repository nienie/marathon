package marathon

import (
    "sync"

    "github.com/nienie/marathon/loadbalancer"
    "github.com/nienie/marathon/httpclient"
)

var (
    cf *clientFactory
)

func init() {
    cf = newClientFactory()
}

type clientFactory struct {
    loadBalancers    map[string]loadbalancer.LoadBalancer
    lbLock           sync.Mutex
    clients          map[string]*httpclient.LoadBalancerHTTPClient
    clientLock       sync.Mutex
}

func newClientFactory() *clientFactory {
    return &clientFactory{
        loadBalancers:      make(map[string]loadbalancer.LoadBalancer),
        lbLock:             sync.Mutex{},
        clients:            make(map[string]*httpclient.LoadBalancerHTTPClient),
        clientLock:         sync.Mutex{},
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
    cf.loadBalancers[name] = lb
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
    cf.clients[name] = client
    cf.clientLock.Unlock()
}

//GetLoadBalancer ...
func GetLoadBalancer(name string) loadbalancer.LoadBalancer {
    return cf.loadBalancers[name]
}

//GetHTTPClient ...
func GetHTTPClient(name string) *httpclient.LoadBalancerHTTPClient {
    return cf.clients[name]
}