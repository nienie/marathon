package stats

import (
    "sync"
)

//Accumulator ...
type Accumulator interface {
    //Publish Called to publish recently collected data.
    Publish()
    //publish is an abstract method, override it~
    publish(*DataBuffer)
    //DataCollector Accumulator must implements DataCollector Interface
    DataCollector
}

//DataAccumulator ...
//DataAccumulator implements DataCollector interface
//used by DataPublisher
type DataAccumulator struct {
    sync.Mutex
    accumulator Accumulator
    //current new data is added to it
    current     *DataBuffer
    //previous used as a source of computed statistics
    previous    *DataBuffer
}

//NewDataAccumulator Actually DataAccumulator is an abstract class.
func NewDataAccumulator(bufferSize int) *DataAccumulator {
    return &DataAccumulator{
        current:    NewDataBuffer(bufferSize),
        previous:   NewDataBuffer(bufferSize),
    }
}

//SetAccumulator ...
func (o *DataAccumulator)SetAccumulator(accumulator Accumulator) {
    o.accumulator = accumulator
}

//NoteValue Accumulating new values
func (o *DataAccumulator)NoteValue(val float64) {
    o.current.NoteValue(val)
}

//Publish Swaps the data collection buffers, and computes statistics about the data collected up til now.
func (o *DataAccumulator)Publish() {
    // lock for prohibit concurrency
    o.Lock()
    //copy current data
    tmpBuffer := o.current.Clone()
    //swap the DataBuffers
    o.current, o.previous = o.previous, o.current
    //start a new collection
    o.current.StartCollection()
    //end the old collection
    tmpBuffer.EndCollection()
    //publish the old collection
    o.accumulator.publish(tmpBuffer)
    //unlock
    o.Unlock()
}

//Dummy implements
func (o *DataAccumulator)publish(*DataBuffer) {}