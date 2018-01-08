package server

import (
	"sync/atomic"
	"time"

	"github.com/nienie/marathon/utils/timer"
)

//PollingServerListUpdater ...
type PollingServerListUpdater struct {
	isActive        int32
	refreshInterval time.Duration
	lastUpdatedTime time.Duration
	refreshTimer    *timer.Timer
}

//NewPollingServerListUpdater ...
func NewPollingServerListUpdater(refreshInterval time.Duration) *PollingServerListUpdater {
	return &PollingServerListUpdater{
		isActive:        int32(0),
		refreshInterval: refreshInterval,
		lastUpdatedTime: time.Duration(0),
	}
}

//Start ...
func (o *PollingServerListUpdater) Start(action UpdateAction) {
	if atomic.CompareAndSwapInt32(&o.isActive, 0, 1) {
		f := func() {
			defer func() {
				if r := recover(); r != nil {

				}
			}()
			action.DoUpdate()
			o.lastUpdatedTime = time.Duration(time.Now().UnixNano())
		}
		o.refreshTimer = timer.NewTimer("PollingServerListUpdater")
		o.refreshTimer.Schedule(f, o.refreshInterval, o.refreshInterval)
	}
}

//Stop ...
func (o *PollingServerListUpdater) Stop() {
	if atomic.CompareAndSwapInt32(&o.isActive, 1, 0) {
		if o.refreshTimer != nil {
			o.refreshTimer.Cancel()
		}
	}
}

//GetLastUpdateTime ...
func (o *PollingServerListUpdater) GetLastUpdateTime() time.Duration {
	return o.lastUpdatedTime
}
