package ratelimit

import (
	"net/url"
	"sync"

	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/server"
)

var (
	leakyBucketLock sync.RWMutex
	leakyBuckets    map[string]*LeakyBucket
)

func init() {
	leakyBuckets = make(map[string]*LeakyBucket)
	leakyBucketLock = sync.RWMutex{}
}

//LeakyBucketRateLimit ...
type LeakyBucketRateLimit struct{}

//NewLeakyBucketRateLimit ...
func NewLeakyBucketRateLimit() *LeakyBucketRateLimit {
	return &LeakyBucketRateLimit{}
}

//Allow ...
func (l *LeakyBucketRateLimit) Allow(uri *url.URL, serverStats *server.Stats, requestConfig config.ClientConfig) bool {
	if uri == nil || serverStats == nil || requestConfig == nil {
		return true
	}

	if !requestConfig.GetPropertyAsBool(config.LeakyBucketRateLimitSwitch, config.DefaultLeakyBucketRateLimitSwitch) {
		return true
	}

	var key string
	if len(uri.Host) == 0 {
		key = serverStats.Server.GetHostPort() + uri.Path
	} else {
		key = uri.Host + uri.Path
	}

	bucket := getLeakyBucket(key, requestConfig)
	if bucket == nil {
		return true
	}

	return bucket.Put()
}

func getLeakyBucket(key string, requestConfig config.ClientConfig) *LeakyBucket {
	bucket, ok := leakyBuckets[key]
	if ok {
		return bucket
	}
	leakyBucketLock.Lock()
	bucket, ok = leakyBuckets[key]
	if ok {
		defer leakyBucketLock.Unlock()
		return bucket
	}
	bucket, err := NewLeakyBucket(
		requestConfig.GetPropertyAsInteger(config.LeakyBucketCapacity, config.DefaultLeakyBucketCapacity),
		requestConfig.GetPropertyAsDuration(config.LeakyBucketInterval, config.DefaultLeakyBucketInterval),
	)
	if err != nil {
		return nil
	}
	leakyBuckets[key] = bucket
	leakyBucketLock.Unlock()
	return bucket
}
