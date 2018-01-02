package stats

import (
	"sync/atomic"
	"time"
)

//DataPublisher ...
type DataPublisher struct {
	dataAccumulator Accumulator
	interval        time.Duration
	stop            chan bool
	isRunning       int32
}

//NewDataPublisher ...
func NewDataPublisher(accumulator Accumulator, interval time.Duration) *DataPublisher {
	return &DataPublisher{
		dataAccumulator: accumulator,
		interval:        interval,
		stop:            make(chan bool),
		isRunning:       int32(0),
	}
}

//Start ...
func (o *DataPublisher) Start() {
	if atomic.CompareAndSwapInt32(&o.isRunning, int32(0), int32(1)) {
		go o.loop()
	}
}

func (o *DataPublisher) loop() {
	defer func() {
		if r := recover(); r != nil {
			atomic.StoreInt32(&o.isRunning, int32(0))
		}
	}()
	ticker := time.NewTicker(time.Duration(o.interval))
	defer ticker.Stop()
	for {
		select {
		case <-o.stop:
			atomic.StoreInt32(&o.isRunning, int32(0))
			return
		case <-ticker.C:
			o.dataAccumulator.Publish()
		}
	}
}

//Stop ...
func (o *DataPublisher) Stop() {
	o.stop <- true
}

//IsRunning ...
func (o *DataPublisher) IsRunning() bool {
	return atomic.LoadInt32(&o.isRunning) > int32(0)
}
