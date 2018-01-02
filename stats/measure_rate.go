package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

//MeasureRate ...
type MeasureRate struct {
	sync.Mutex
	lastBucket     int64
	currentBucket  int64
	sampleInterval time.Duration
	threshold      time.Duration
}

//NewMeasureRate sampleInterval in milliseconds
func NewMeasureRate(sampleInterval time.Duration) *MeasureRate {
	return &MeasureRate{
		lastBucket:     int64(0),
		currentBucket:  int64(0),
		sampleInterval: sampleInterval,
		threshold:      time.Duration(time.Now().UnixNano()) + sampleInterval,
	}
}

//GetCount Returns the count in the last sample interval
func (o *MeasureRate) GetCount() int64 {
	o.checkAndResetWindow()
	return atomic.LoadInt64(&o.lastBucket)
}

func (o *MeasureRate) checkAndResetWindow() {
	now := time.Duration(time.Now().UnixNano())
	o.Lock()
	defer o.Unlock()
	if o.threshold < now {
		atomic.SwapInt64(&o.lastBucket, o.currentBucket)
		atomic.StoreInt64(&o.currentBucket, int64(0))
		o.threshold = now + o.sampleInterval
	}
}

//GetCurrentCount Returns the count in the current sample interval which will be incomplete.
func (o *MeasureRate) GetCurrentCount() int64 {
	o.checkAndResetWindow()
	return atomic.LoadInt64(&o.currentBucket)
}

//Increment Increments the count in the current sample interval.
func (o *MeasureRate) Increment() {
	o.checkAndResetWindow()
	atomic.AddInt64(&o.currentBucket, int64(1))
}
