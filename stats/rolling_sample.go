
package stats

import (
    "sync"
    "time"

    "github.com/rcrowley/go-metrics"
)

//RollingSample ...
type RollingSample struct {
    sync.RWMutex
    Buckets     map[int64][]int64   //key is a unix timestamp(unit: seconds)
    WindowSize  int                 //save for how long time...(unit: seconds)
}

//NewRollingSample ...
func NewRollingSample(windowSize int) *RollingSample {
    return &RollingSample{
        Buckets:        make(map[int64][]int64),
        WindowSize:     windowSize,
    }
}

//UpdateValue ...
func (s *RollingSample)UpdateValue(val int64) {
    currentTime := time.Now().Unix()
    s.Lock()
    defer s.Unlock()
    _, ok := s.Buckets[currentTime]
    if !ok {
        s.Buckets[currentTime] = make([]int64, 0, 1000)
    }
    s.Buckets[currentTime] = append(s.Buckets[currentTime], val)
    s.removeExpiredBuckets(currentTime)
}

func (s *RollingSample)removeExpiredBuckets(currentTime int64) {
    if len(s.Buckets) < s.WindowSize {
        return
    }
    for timestamp := range s.Buckets {
        if currentTime - timestamp > int64(s.WindowSize) {
            delete(s.Buckets, timestamp)
        }
    }
}

func (s *RollingSample)getHistoryData() []int64 {
    values := make([]int64, 0, 1000)
    s.RLock()
    for _, data := range s.Buckets {
        values = append(values, data...)
    }
    s.RUnlock()
    return values
}

//Max ...
func (s *RollingSample)Max() int64 {
    return metrics.SampleMax(s.getHistoryData())
}

//Min ...
func (s *RollingSample)Min() int64 {
    return metrics.SampleMin(s.getHistoryData())
}

//Mean ...
func (s *RollingSample)Mean() float64 {
    return metrics.SampleMean(s.getHistoryData())
}

//StdDev ...
func (s *RollingSample)StdDev() float64 {
    return metrics.SampleStdDev(s.getHistoryData())
}

//Variance ...
func (s *RollingSample)Variance() float64 {
    return metrics.SampleVariance(s.getHistoryData())
}

//Percentile ...
func (s *RollingSample)Percentile(p float64) float64 {
    return metrics.SamplePercentile(s.getHistoryData(), p)
}

//Percentiles ...
func (s *RollingSample)Percentiles(ps []float64) []float64 {
    return metrics.SamplePercentiles(s.getHistoryData(), ps)
}