package ratelimit

import (
	"fmt"
	"time"
)

//LeakyBucket ...
type LeakyBucket struct {
	capacity int
	ch       chan struct{}
	stop     chan bool
	interval time.Duration
}

//NewLeakyBucket ...
func NewLeakyBucket(capacity int, interval time.Duration) (*LeakyBucket, error) {
	if capacity <= 0 || interval <= time.Duration(0) {
		return nil, fmt.Errorf("invalid parameters")
	}

	bucket := &LeakyBucket{
		capacity: capacity,
		ch:       make(chan struct{}, capacity),
		stop:     make(chan bool, 1),
		interval: interval,
	}
	go bucket.run()
	return bucket, nil
}

func (b *LeakyBucket) run() {
	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-b.stop:
			return
		case <-t.C:
			select {
			case <-b.ch:
			default:
			}
		}
	}
}

//Put ...
func (b *LeakyBucket) Put() bool {
	select {
	case b.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

//Stop ...
func (b *LeakyBucket) Stop() {
	b.stop <- true
}
