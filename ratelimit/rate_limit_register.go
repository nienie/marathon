package ratelimit

import (
    "net/url"

    "github.com/nienie/marathon/config"
    "github.com/nienie/marathon/server"
)
var (
    rateLimitRegister []RateLimit
)

func init() {
    rateLimitRegister = make([]RateLimit, 0)
    rateLimitRegister = append(rateLimitRegister, NewConcurrencyRateLimit())
    rateLimitRegister = append(rateLimitRegister, NewTokenBucketRateLimit())
    rateLimitRegister = append(rateLimitRegister, NewLeakyBucketRateLimit())
}

//Allow ...
func Allow(url *url.URL, serverStats *server.Stats, requestConfig config.ClientConfig) bool {
    for _, rateLimit := range rateLimitRegister {
        if rateLimit.Allow(url, serverStats, requestConfig) == false {
            return false
        }
    }
    return true
}

//RegisterRateLimit ...
func RegisterRateLimit(r RateLimit) {
    rateLimitRegister = append(rateLimitRegister, r)
}
