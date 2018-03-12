package ratelimit

import (
	"fmt"
	"time"
)

const (
	//DefaultFillInterval ...
	DefaultFillInterval = 10 * time.Millisecond
)

//TokenBucket ...
type TokenBucket struct {
	ch           chan struct{}
	capacity     int
	fillCount    int
	fillInterval time.Duration
	stop         chan bool
}

//NewTokenBucket ...
func NewTokenBucket(capacity int, fillInterval time.Duration, fillCount int) (*TokenBucket, error) {
	if capacity <= 0 || fillCount <= 0 {
		return nil, fmt.Errorf("invalid parameters")
	}

	bucket := &TokenBucket{
		capacity:     capacity,
		fillInterval: fillInterval,
		fillCount:    fillCount,
		ch:           make(chan struct{}, capacity),
		stop:         make(chan bool, 1),
	}
	for i := 0; i < fillCount; i++ {
		bucket.ch <- struct{}{}
	}
	go bucket.run()
	return bucket, nil
}

func (b *TokenBucket) run() {
	t := time.NewTicker(b.fillInterval)
	defer t.Stop()
	for {
		select {
		case <-b.stop:
			return
		case <-t.C:
			for i := 0; i < b.fillCount; i++ {
				select {
				case b.ch <- struct{}{}:
				default:
				}
			}
		}
	}
}

//GetToken ...
func (b *TokenBucket) GetToken() bool {
	select {
	case <-b.ch:
		return true
	default:
		return false
	}
}

//Stop ...
func (b *TokenBucket) Stop() {
	b.stop <- true
}
