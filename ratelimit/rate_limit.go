package ratelimit

import (
    "net/url"

    "github.com/nienie/marathon/config"
    "github.com/nienie/marathon/server"
)

//RateLimit ...
type RateLimit interface {

    //Allow ...
    Allow(url *url.URL, serverStats *server.Stats, requestConfig config.ClientConfig) bool
}