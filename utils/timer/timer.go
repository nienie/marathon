package timer

import (
	"runtime"
	"sync/atomic"
	"time"
)

//Timer ...
type Timer struct {
	Name      string
	interval  time.Duration
	stop      chan bool
	isRunning int32
}

//NewTimer ...
func NewTimer(name string) *Timer {
	return &Timer{
		Name:      name,
		stop:      make(chan bool, 1),
		isRunning: int32(0),
	}
}

//Schedule ...
func (t *Timer) Schedule(task func(), period time.Duration, delay time.Duration) {
	if task == nil || period <= time.Duration(0) || delay < time.Duration(0) {
		return
	}

	if atomic.CompareAndSwapInt32(&t.isRunning, int32(0), int32(1)) {
		go func() {
			if delay > 0 {
				dt := time.After(delay)
				for {
					select {
					case <-t.stop:
						return
					case <-dt:
					}
				}
			}
			ticker := time.NewTicker(period)
			defer ticker.Stop()
			runtime.SetFinalizer(t, func(t *Timer) {
				t.Cancel()
			})
			for {
				select {
				case <-t.stop:
					return
				case <-ticker.C:
					go task()
				}
			}
		}()
	}
}

//Cancel ...
func (t *Timer) Cancel() {
	if atomic.CompareAndSwapInt32(&t.isRunning, int32(1), int32(0)) {
		//stop the Timer, does not close the channel, to prevent a read from the channel succeeding incorrectly.
		t.stop <- true
	}
}
