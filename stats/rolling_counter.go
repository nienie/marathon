package stats

import (
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

//RollingCounter ...
type RollingCounter struct {
	*sync.RWMutex
	Buckets    map[int64]metrics.Counter //key is a unix timestamp(unit: seconds)
	WindowSize int                       //save for how long time...(unit: seconds)
}

//NewRollingCounter ...
func NewRollingCounter(windowSize int) *RollingCounter {
	return &RollingCounter{
		RWMutex:    &sync.RWMutex{},
		Buckets:    make(map[int64]metrics.Counter),
		WindowSize: windowSize,
	}
}

//Inc ...
func (c *RollingCounter) Inc(i int64) {
	currentTime := time.Now().Unix()
	c.Lock()
	defer c.Unlock()

	//update counter
	_, ok := c.Buckets[currentTime]
	if !ok {
		c.Buckets[currentTime] = &metrics.StandardCounter{}
	}
	c.Buckets[currentTime].Inc(i)

	//delete the expired counter
	c.removeExpiredBuckets(currentTime)
}

func (c *RollingCounter) removeExpiredBuckets(currentTime int64) {
	//delete the expired counter
	if len(c.Buckets) < c.WindowSize {
		return
	}
	for timestamp := range c.Buckets {
		if currentTime-timestamp > int64(c.WindowSize) {
			delete(c.Buckets, timestamp)
		}
	}
}

//Dec ...
func (c *RollingCounter) Dec(i int64) {
	currentTime := time.Now().Unix()
	c.Lock()
	defer c.Unlock()

	//update and counter
	_, ok := c.Buckets[currentTime]
	if !ok {
		c.Buckets[currentTime] = &metrics.StandardCounter{}
	}
	c.Buckets[currentTime].Dec(i)

	//delete the expired counter
	c.removeExpiredBuckets(currentTime)
}

//Sum ...
func (c *RollingCounter) Sum(size int) int64 {
	if size > c.WindowSize || size < 0 {
		size = c.WindowSize
	}
	var totalCount = int64(0)
	currentTime := time.Now().Unix()

	c.RLock()
	defer c.RUnlock()
	for timestamp, counter := range c.Buckets {
		if currentTime-timestamp <= int64(size) {
			totalCount += counter.Count()
		}
	}
	return totalCount
}

//Clear ...
func (c *RollingCounter) Clear() {
	c.Lock()
	defer c.Unlock()
	c.Buckets = make(map[int64]metrics.Counter)
}

//Count ...
func (c *RollingCounter) Count() int64 {
	return c.Sum(c.WindowSize)
}

//Snapshot ...
func (c *RollingCounter) Snapshot() metrics.Counter {
	cc := NewRollingCounter(c.WindowSize)
	c.RLock()
	for timestamp, counter := range c.Buckets {
		cc.Buckets[timestamp] = counter.Snapshot()
	}
	return cc
}
