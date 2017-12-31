package stats

import (
    "sync"
    "time"
)

//DataDistribution publishes statistics about the "previous" buffer
//DataDistribution extends DataAccumulator
type DataDistribution struct {
    sync.Mutex
    *DataAccumulator
    NumValues   int64
    Mean        float64
    Variance    float64
    Stddev      float64
    Min         float64
    Max         float64
    Timestamp   time.Duration
    Percents    []float64
    Percentiles []float64

    interval    time.Duration
    size        int
}

//NewDataDistribution ...
func NewDataDistribution(bufferSize int, percents []float64) *DataDistribution {
    dataDistribution := &DataDistribution{}
    dataDistribution.DataAccumulator = NewDataAccumulator(bufferSize)
    dataDistribution.DataAccumulator.SetAccumulator(dataDistribution)
    dataDistribution.Percents = percents
    dataDistribution.Percentiles = make([]float64, len(percents))
    return dataDistribution
}

// override Accumulator publish method
func (o *DataDistribution)publish(buf *DataBuffer) {
    o.Timestamp = time.Duration(time.Now().UnixNano())
    o.NumValues = buf.GetNumValues()
    o.Mean = buf.GetMean()
    o.Variance = buf.GetVariance()
    o.Stddev = buf.GetStdDev()
    o.Min = buf.GetMinimum()
    o.Max = buf.GetMaximum()
    o.Percentiles = buf.GetPercentiles(o.Percents)
    o.interval = buf.GetSampleInterval()
    o.size = buf.GetSampleSize()
}

//Clear ...
func (o *DataDistribution)Clear() {
    o.NumValues = int64(0)
    o.Mean =float64(0.0)
    o.Variance = float64(0.0)
    o.Stddev = float64(0.0)
    o.Min = float64(0.0)
    o.Max = float64(0.0)
    o.Timestamp = time.Duration(0)
    for i := 0; i < len(o.Percentiles); i++ {
        o.Percentiles[i] = float64(0.0)
    }
    o.interval = time.Duration(0)
    o.size = 0
}

//GetSampleInterval ...
func (o *DataDistribution)GetSampleInterval() time.Duration {
    return o.interval
}

//GetSampleSize ...
func (o *DataDistribution)GetSampleSize() int {
    return o.size
}

//GetNumValues ...
func (o *DataDistribution)GetNumValues() int64 {
    return o.NumValues
}

//GetMean ...
func (o *DataDistribution)GetMean() float64 {
    return o.Mean
}

//GetVariance ...
func (o *DataDistribution)GetVariance() float64 {
    return o.Variance
}

//GetStdDev ...
func (o *DataDistribution)GetStdDev() float64 {
    return o.Stddev
}

//GetMinimum ...
func (o *DataDistribution)GetMinimum() float64 {
    return o.Min
}

//GetMaximum ...
func (o *DataDistribution)GetMaximum() float64 {
    return o.Max
}

//GetTimestamp ...
func (o *DataDistribution)GetTimestamp() time.Duration {
    return o.Timestamp
}

//GetPercents ...
func (o *DataDistribution)GetPercents() []float64 {
    return o.Percents
}

//GetPercentiles ...
func (o *DataDistribution)GetPercentiles() []float64 {
    return o.Percentiles
}