package ratelimit

import (
	"net/url"
	"sync"

	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/server"
)

var (
	tokenBucketLock sync.RWMutex
	tokenBuckets    map[string]*TokenBucket
)

func init() {
	tokenBuckets = make(map[string]*TokenBucket)
	tokenBucketLock = sync.RWMutex{}
}

//TokenBucketRateLimit ...
type TokenBucketRateLimit struct{}

//NewTokenBucketRateLimit ...
func NewTokenBucketRateLimit() *TokenBucketRateLimit {
	return &TokenBucketRateLimit{}
}

//Allow ...
func (l *TokenBucketRateLimit) Allow(uri *url.URL, serverStats *server.Stats, requestConfig config.ClientConfig) bool {
	if uri == nil || requestConfig == nil {
		return true
	}

	if !requestConfig.GetPropertyAsBool(config.TokenBucketRateLimitSwitch, config.DefaultTokenBucketRateLimitSwitch) {
		return true
	}

	var key = uri.Path
	//if len(uri.Host) == 0 {
	//	key = serverStats.Server.GetHostPort() + uri.Path
	//} else {
	//	key = uri.Host + uri.Path
	//}

	bucket := getTokenBucket(key, requestConfig)
	if bucket == nil {
		return true
	}

	return bucket.GetToken()
}

func getTokenBucket(key string, requestConfig config.ClientConfig) *TokenBucket {
	bucket, ok := tokenBuckets[key]
	if ok {
		return bucket
	}
	tokenBucketLock.Lock()
	bucket, ok = tokenBuckets[key]
	if ok {
		defer tokenBucketLock.Unlock()
		return bucket
	}
	bucket, err := NewTokenBucket(
		requestConfig.GetPropertyAsInteger(config.TokenBucketCapacity, config.DefaultTokenBucketCapacity),
		requestConfig.GetPropertyAsDuration(config.TokenBucketFillInterval, config.DefaultTokenBucketFillInterval),
		requestConfig.GetPropertyAsInteger(config.TokenBucketFillCount, config.DefaultTokenBucketFillCount),
	)
	if err != nil {
		return nil
	}
	tokenBuckets[key] = bucket
	tokenBucketLock.Unlock()
	return bucket
}
