package monitor

import "time"

//StopWatch measures the time token for execution of some code.
type StopWatch interface {
	//Start mark the start time
	Start()

	//Stop mark the end time
	Stop()

	//reset the stopwatch so that it can be used again
	Reset()

	//GetDuration
	GetDuration() time.Duration
}

//BasicStopWatch ...
type BasicStopWatch struct {
	startTime time.Duration
	endTime   time.Duration
	running   bool
}

//NewBasicStopWatch ...
func NewBasicStopWatch() StopWatch {
	return &BasicStopWatch{}
}

//Start ...
func (w *BasicStopWatch) Start() {
	w.startTime = time.Duration(time.Now().UnixNano())
	w.running = true
}

//Stop ...
func (w *BasicStopWatch) Stop() {
	w.endTime = time.Duration(time.Now().UnixNano())
	w.running = false
}

//Reset ...
func (w *BasicStopWatch) Reset() {
	w.startTime = time.Duration(0)
	w.endTime = time.Duration(0)
	w.running = false
}

//GetDuration ...
func (w *BasicStopWatch) GetDuration() time.Duration {
	end := time.Duration(time.Now().UnixNano())
	if !w.running {
		end = w.endTime
	}
	return end - w.startTime
}
