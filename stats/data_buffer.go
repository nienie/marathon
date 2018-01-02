package stats

import (
	"math"
	"sort"
	"sync"
	"time"
)

//DataBuffer A fixed-size data collection buffer that holds a sliding window of the most recent values added.
// The struct implements DataCollector and extends Distribution.
type DataBuffer struct {
	sync.Mutex
	*Distribution
	startTime time.Duration
	endTime   time.Duration
	size      int
	capacity  int
	buffer    []float64
	insertPos int
}

//NewDataBuffer ...
func NewDataBuffer(capacity int) *DataBuffer {
	return &DataBuffer{
		Distribution: NewDistribution(),
		startTime:    time.Duration(0),
		endTime:      time.Duration(0),
		size:         0,
		capacity:     capacity,
		buffer:       make([]float64, capacity),
		insertPos:    0,
	}
}

//GetSampleInterval Gets the length of time over which the data was collected in nanoseconds
//The value is only valid after EndCollection, has been called (and before a subsequent call to StartCollection).
func (o *DataBuffer) GetSampleInterval() time.Duration {
	return o.endTime - o.startTime
}

//Clear Override Distribution Clear()
func (o *DataBuffer) Clear() {
	o.Distribution.Clear()
	o.startTime = time.Duration(0)
	o.endTime = time.Duration(0)
	o.size = 0
	o.insertPos = 0
}

//GetCapacity ...
func (o *DataBuffer) GetCapacity() int {
	return o.capacity
}

//GetSampleSize Gets the number of values currently held in the buffer.
func (o *DataBuffer) GetSampleSize() int {
	return o.size
}

//StartCollection Notifies the buffer that data is collection is now enabled.
func (o *DataBuffer) StartCollection() {
	o.Lock()
	defer o.Unlock()
	o.Clear()
	o.startTime = time.Duration(time.Now().UnixNano())
}

//for sort
type sortInterface []float64

func (s sortInterface) Len() int {
	return len(s)
}
func (s sortInterface) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortInterface) Less(i, j int) bool {
	return s[i] < s[j]
}

//EndCollection Notifies the buffer that data has just ended.
func (o *DataBuffer) EndCollection() {
	o.endTime = time.Duration(time.Now().UnixNano())
	buf := sortInterface(o.buffer[:o.size])
	sort.Sort(buf)
	copy(o.buffer, []float64(buf))
}

//NoteValue The buffer wraps-around if it is full, overwriting the oldest entry with the new value.
//Override struct Distribution NoteValue function
func (o *DataBuffer) NoteValue(val float64) {
	o.Lock()
	defer o.Unlock()

	o.Distribution.NoteValue(val)
	o.buffer[o.insertPos] = val
	o.insertPos++
	if o.insertPos >= o.capacity {
		o.insertPos = 0
		o.size = o.capacity
	} else if o.insertPos > o.size {
		o.size = o.insertPos
	}
}

//GetPercentiles Gets the requested percentile statistics.
func (o *DataBuffer) GetPercentiles(percents []float64) []float64 {
	percentiles := make([]float64, len(percents))
	for i := 0; i < len(percents); i++ {
		percentiles[i] = o.computePercentile(percents[i])
	}
	return percentiles
}

func (o *DataBuffer) computePercentile(percent float64) float64 {
	if o.size <= 0 {
		return float64(0.0)
	} else if percent <= float64(0.0) {
		return o.buffer[0]
	} else if percent >= float64(100.0) {
		return o.buffer[o.size-1]
	}

	index := (percent / float64(100.0)) * float64(o.size)
	iLow := int(math.Floor(index))
	iHigh := int(math.Ceil(index))
	if iHigh > o.size {
		return o.buffer[o.size-1]
	} else if iLow == iHigh {
		return o.buffer[iLow]
	} else {
		// Interpolate between the two bounding values
		return o.buffer[iLow] + (index-float64(iLow))*(o.buffer[iHigh]-o.buffer[iLow])
	}
}

//Clone clone self
func (o *DataBuffer) Clone() *DataBuffer {
	dataBuffer := &DataBuffer{}
	dataBuffer.Distribution = o.Distribution.Clone()
	dataBuffer.startTime = o.startTime
	dataBuffer.endTime = o.endTime
	dataBuffer.size = o.size
	dataBuffer.capacity = o.capacity
	dataBuffer.buffer = make([]float64, o.capacity)
	copy(dataBuffer.buffer, o.buffer)
	dataBuffer.insertPos = o.insertPos
	return dataBuffer
}
