package stats

import (
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"sort"
)

//RollingSample ...
type RollingSample struct {
	*sync.RWMutex
	Buckets    map[int64][]int64 //key is a unix timestamp(unit: seconds)
	WindowSize int               //save for how long time...(unit: seconds)
}

//NewRollingSample ...
func NewRollingSample(windowSize int) *RollingSample {
	return &RollingSample{
		RWMutex:    &sync.RWMutex{},
		Buckets:    make(map[int64][]int64),
		WindowSize: windowSize,
	}
}

//UpdateValue ...
func (s *RollingSample) UpdateValue(val int64) {
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

func (s *RollingSample) removeExpiredBuckets(currentTime int64) {
	if len(s.Buckets) < s.WindowSize {
		return
	}
	for timestamp := range s.Buckets {
		if currentTime-timestamp > int64(s.WindowSize) {
			delete(s.Buckets, timestamp)
		}
	}
}

func (s *RollingSample) getHistoryData() []int64 {
	values := make([]int64, 0, 1000)
	currentTime := time.Now().Unix()
	s.RLock()
	for timestamp, data := range s.Buckets {
		if currentTime-timestamp > int64(s.WindowSize) {
			continue
		}
		values = append(values, data...)
	}
	s.RUnlock()
	return values
}

//Sum ...
func (s *RollingSample) Sum() int64 {
	return metrics.SampleSum(s.getHistoryData())
}

//Max ...
func (s *RollingSample) Max() int64 {
	return metrics.SampleMax(s.getHistoryData())
}

//Min ...
func (s *RollingSample) Min() int64 {
	return metrics.SampleMin(s.getHistoryData())
}

//Mean ...
func (s *RollingSample) Mean() float64 {
	return metrics.SampleMean(s.getHistoryData())
}

//StdDev ...
func (s *RollingSample) StdDev() float64 {
	return metrics.SampleStdDev(s.getHistoryData())
}

//Variance ...
func (s *RollingSample) Variance() float64 {
	return metrics.SampleVariance(s.getHistoryData())
}

//Percentile ...
func (s *RollingSample) Percentile(p float64) float64 {
	return metrics.SamplePercentile(s.getHistoryData(), p)
}

//Percentiles ...
func (s *RollingSample) Percentiles(ps []float64) []float64 {
	return metrics.SamplePercentiles(s.getHistoryData(), ps)
}

//AvgPerSecond ...
func (s *RollingSample) AvgPerSecond() []float64 {
	currentTime := time.Now().Unix()
	ret := make([]float64, 0, s.WindowSize)
	var timestamps byInt64
	s.RLock()
	for timestamp := range s.Buckets {
		if currentTime-timestamp > int64(s.WindowSize) {
			continue
		}
		timestamps = append(timestamps, timestamp)
	}
	sort.Sort(timestamps)
	for _, timestamp := range timestamps {
		ret = append(ret, metrics.SampleMean(s.Buckets[timestamp]))
	}
	s.RUnlock()
	return ret
}

type byInt64 []int64

func (c byInt64) Len() int           { return len(c) }
func (c byInt64) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byInt64) Less(i, j int) bool { return c[i] < c[j] }