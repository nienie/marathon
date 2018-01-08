package ratelimit

import (
    "net/url"

    "github.com/nienie/marathon/config"
    "github.com/nienie/marathon/server"
)

//ConcurrencyRateLimit ...
type ConcurrencyRateLimit struct {}

//NewConcurrencyRateLimit ...
func NewConcurrencyRateLimit() *ConcurrencyRateLimit {
    return &ConcurrencyRateLimit{}
}

//Allow ...
func (l *ConcurrencyRateLimit)Allow(url *url.URL, serverStats *server.Stats, requestConfig config.ClientConfig) bool{
    if serverStats == nil || requestConfig == nil {
        return true
    }

    if requestConfig.GetPropertyAsBool(config.ConcurrencyRateLimitSwitch, config.DefaultConcurrencyRateLimitSwitch) == false{
        return true
    }

    maxConns := int64(requestConfig.GetPropertyAsInteger(config.MaxConnectionsPerHost, config.DefaultMaxConnectionsPerHost))
    curConns := serverStats.GetOpenConnectionsCount()
    if curConns > maxConns {
        return false
    }
    return true
}
