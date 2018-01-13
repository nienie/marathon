package server

import (
	"sync/atomic"
	"time"

	"github.com/nienie/marathon/stats"
	"github.com/rcrowley/go-metrics"
)

const (
	tenPercentile                 float64 = 10
	twentyFivePercentile          float64 = 25
	fiftyPercentile               float64 = 50
	seventyFivePercentile         float64 = 75
	ninetyPercentile              float64 = 90
	ninetyFivePercentile          float64 = 95
	ninetyEightPercentile         float64 = 98
	ninetyNinePercentile          float64 = 99
	ninetyNinePointFivePercentile float64 = 99.5
)

const (
	//DefaultConnectionFailureCountThreshold ...
	DefaultConnectionFailureCountThreshold = 5
	//DefaultCircuitTrippedTimeoutFactor ...
	DefaultCircuitTrippedTimeoutFactor = 10
	//DefaultMaxCircuitTrippedTimeout ...
	DefaultMaxCircuitTrippedTimeout = 10 * time.Second
	//DefaultActiveRequestsCountTimeout ...
	DefaultActiveRequestsCountTimeout = 30 * time.Second
	//DefaultRequestCountsSlidingWindowSize ...
	DefaultRequestCountsSlidingWindowSize = 300 //store 300 seconds' data
	//DefaultResponseTimeWindowSize ...
	DefaultResponseTimeWindowSize = 300 //store 300 seconds' data
)

//Stats ...
type Stats struct {
	Server *Server

	ConnectionFailureThreshold  int
	CircuitTrippedTimeoutFactor int
	MaxCircuitTrippedTimeout    time.Duration
	ActiveRequestsCountTimeout  time.Duration

	RequestCountsSlidingWindowSize int
	ResponseTimeWindowSize 		   int


	//for stats
	totalRequests                     metrics.Counter
	activeRequestsCount               metrics.Counter
	openConnectionsCount              metrics.Counter
	successiveConnectionFailureCount  metrics.Counter
	totalCircuitBreakerBlackOutPeriod int64 //nanoseconds

	//record time
	lastConnectionFailedTimestamp          int64
	lastActiveRequestsCountChangeTimestamp int64
	firstConnectionTimestamp               int64
	lastAccessedTimestamp                  int64

	//stats objects
	responseTimeDist     *stats.Distribution  //to stats in the overall time
	responseTimeInWindow *stats.RollingSample // to stats in a recent time-slice

	serverFailureCounts  *stats.RollingCounter //server failure counts in a sliding window time
	requestCountInWindow *stats.RollingCounter //request count in a window time
}

//NewDefaultServerStats ...
func NewDefaultServerStats() *Stats {
	return &Stats{
		ConnectionFailureThreshold:  DefaultConnectionFailureCountThreshold,
		CircuitTrippedTimeoutFactor: DefaultCircuitTrippedTimeoutFactor,
		MaxCircuitTrippedTimeout:    DefaultMaxCircuitTrippedTimeout,
		ActiveRequestsCountTimeout:  DefaultActiveRequestsCountTimeout,

		RequestCountsSlidingWindowSize: DefaultRequestCountsSlidingWindowSize,
		ResponseTimeWindowSize:			DefaultResponseTimeWindowSize,

		responseTimeDist:            	  stats.NewDistribution(),

		totalRequests:                    &metrics.StandardCounter{},
		activeRequestsCount:              &metrics.StandardCounter{},
		openConnectionsCount:             &metrics.StandardCounter{},
		successiveConnectionFailureCount: &metrics.StandardCounter{},
	}
}

//Initialize ...
func (o *Stats) Initialize(svr *Server) {
	o.Server = svr

	o.serverFailureCounts =  stats.NewRollingCounter(o.RequestCountsSlidingWindowSize)
	o.requestCountInWindow = stats.NewRollingCounter(o.RequestCountsSlidingWindowSize)
	o.responseTimeInWindow = stats.NewRollingSample(o.ResponseTimeWindowSize)
}

//Close ...
func (o *Stats) Close() {}

//AddToFailureCount increment the count of failure for this server
func (o *Stats) AddToFailureCount() {
	o.serverFailureCounts.Inc(int64(1))
}

//GetFailureCount returns the count of failures in the current window.
func (o *Stats) GetFailureCount() int64 {
	return o.serverFailureCounts.Count()
}

//NoteResponseTime call this method to note the response time after every request.
func (o *Stats) NoteResponseTime(msecs float64) {
	o.responseTimeDist.NoteValue(msecs)
	o.responseTimeInWindow.UpdateValue(int64(msecs))
}

//IncrementNumRequests note the total number of requests.
func (o *Stats) IncrementNumRequests() {
	o.totalRequests.Inc(int64(1))
}

//IncrementActiveRequestsCount note the active number of requests.
func (o *Stats) IncrementActiveRequestsCount() {
	o.activeRequestsCount.Inc(int64(1))
	o.requestCountInWindow.Inc(int64(1))
	currentTime := time.Now().UnixNano()
	atomic.StoreInt64(&o.lastActiveRequestsCountChangeTimestamp, currentTime)
	atomic.StoreInt64(&o.lastAccessedTimestamp, currentTime)
	if o.firstConnectionTimestamp == int64(0) {
		atomic.StoreInt64(&o.firstConnectionTimestamp, currentTime)
	}
}

//DecrementActiveRequestsCount ...
func (o *Stats) DecrementActiveRequestsCount() {
	o.activeRequestsCount.Dec(int64(1))
	if o.activeRequestsCount.Count() < int64(0) {
		o.activeRequestsCount.Clear()
	}
	atomic.StoreInt64(&o.lastActiveRequestsCountChangeTimestamp, time.Now().UnixNano())
}

//GetActiveRequestsCount ...
func (o *Stats) GetActiveRequestsCount(currentTime time.Duration) int64 {
	count := o.activeRequestsCount.Count()

	if currentTime-time.Duration(o.lastActiveRequestsCountChangeTimestamp) > o.ActiveRequestsCountTimeout || count < 0 {
		o.activeRequestsCount.Clear()
		return 0
	}

	return count
}

//IncrementOpenConnectionsCount ...
func (o *Stats) IncrementOpenConnectionsCount() {
	o.openConnectionsCount.Inc(int64(1))
}

//DecrementOpenConnectionsCount ...
func (o *Stats) DecrementOpenConnectionsCount() {
	o.openConnectionsCount.Dec(int64(1))
	if o.openConnectionsCount.Count() < 0 {
		o.openConnectionsCount.Clear()
	}
}

//GetOpenConnectionsCount ...
func (o *Stats) GetOpenConnectionsCount() int64 {
	return o.openConnectionsCount.Count()
}

//GetMeasuredRequestsCount ...
func (o *Stats) GetMeasuredRequestsCount() int64 {
	return o.requestCountInWindow.Count()
}

func (o *Stats) getCircuitBreakerBlackoutPeriod() time.Duration {
	failureCount := o.successiveConnectionFailureCount.Count()
	if failureCount < int64(o.ConnectionFailureThreshold) {
		return time.Duration(0)
	}

	diff := failureCount - int64(o.ConnectionFailureThreshold)
	if diff > 16 {
		diff = 16
	}

	blackOutSeconds := time.Duration(int64(o.CircuitTrippedTimeoutFactor)*diff*2) * time.Second
	if blackOutSeconds > o.MaxCircuitTrippedTimeout {
		blackOutSeconds = o.MaxCircuitTrippedTimeout
	}

	return blackOutSeconds
}

func (o *Stats) getCircuitBreakerTimeout() time.Duration {
	blackOutPeriod := o.getCircuitBreakerBlackoutPeriod()

	if blackOutPeriod <= time.Duration(0) {
		return time.Duration(0)
	}

	return time.Duration(o.lastConnectionFailedTimestamp) + blackOutPeriod
}

//IsCircuitBreakerTripped ...
func (o *Stats) IsCircuitBreakerTripped(currentTime time.Duration) bool {
	circuitBreakerTimeout := o.getCircuitBreakerTimeout()
	if circuitBreakerTimeout <= time.Duration(0) {
		return false
	}
	return circuitBreakerTimeout > currentTime
}

//IncrementSuccessiveConnectionFailureCount ...
func (o *Stats) IncrementSuccessiveConnectionFailureCount() {
	atomic.StoreInt64(&o.lastConnectionFailedTimestamp, time.Now().UnixNano())
	o.successiveConnectionFailureCount.Inc(int64(1))
	atomic.AddInt64(&o.totalCircuitBreakerBlackOutPeriod, int64(o.getCircuitBreakerBlackoutPeriod()))
}

//ClearSuccessiveConnectionFailureCount ...
func (o *Stats) ClearSuccessiveConnectionFailureCount() {
	o.successiveConnectionFailureCount.Clear()
}

//GetSuccessiveConnectionCount ...
func (o *Stats) GetSuccessiveConnectionCount() int64 {
	return o.successiveConnectionFailureCount.Count()
}

//GetResponseTimeAvg gets the average total amount of time to handle a request, in milliseconds.
func (o *Stats) GetResponseTimeAvg() float64 {
	return o.responseTimeDist.GetMean()
}

//GetResponseTimeMax gets the maximum amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTimeMax() float64 {
	return o.responseTimeDist.GetMaximum()
}

//GetResponseTimeMin gets the minimum amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTimeMin() float64 {
	return o.responseTimeDist.GetMinimum()
}

//GetResponseTimeStdDev gets the standard deviation in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTimeStdDev() float64 {
	return o.responseTimeDist.GetStdDev()
}

//GetResponseTimeAvgRecent gets the average total amount of time to handle a request in the recent time-slice, in milliseconds.
func (o *Stats) GetResponseTimeAvgRecent() float64 {
	return o.responseTimeInWindow.Mean()
}

func (o *Stats) getResponseTimePercentile(percent float64) float64 {
	return o.responseTimeInWindow.Percentile(percent / float64(100))
}

//GetResponseTime10thPercentile gets the 10-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime10thPercentile() float64 {
	return o.getResponseTimePercentile(tenPercentile)
}

//GetResponseTime25thPercentile gets the 25-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime25thPercentile() float64 {
	return o.getResponseTimePercentile(twentyFivePercentile)
}

//GetResponseTime50thPercentile gets the 50-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime50thPercentile() float64 {
	return o.getResponseTimePercentile(fiftyPercentile)
}

//GetResponseTime75thPercentile gets the 75-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime75thPercentile() float64 {
	return o.getResponseTimePercentile(seventyFivePercentile)
}

//GetResponseTime90thPercentile gets the 90-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime90thPercentile() float64 {
	return o.getResponseTimePercentile(ninetyPercentile)
}

//GetResponseTime95thPercentile gets the 95-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime95thPercentile() float64 {
	return o.getResponseTimePercentile(ninetyFivePercentile)
}

//GetResponseTime98thPercentile gets the 98-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime98thPercentile() float64 {
	return o.getResponseTimePercentile(ninetyEightPercentile)
}

//GetResponseTime99thPercentile gets the 99-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime99thPercentile() float64 {
	return o.getResponseTimePercentile(ninetyNinePercentile)
}

//GetResponseTime99point5thPercentile gets the 99.5-th percentile in the total amount of time spent handling a request, in milliseconds.
func (o *Stats) GetResponseTime99point5thPercentile() float64 {
	return o.getResponseTimePercentile(ninetyNinePointFivePercentile)
}

//GetTotalRequestsCount ...
func (o *Stats) GetTotalRequestsCount() int64 {
	return o.totalRequests.Count()
}

//GetErrorPercentage ...
func (o *Stats)GetErrorPercentage(size int) float64 {
	errorCount := o.serverFailureCounts.Sum(size)
	totalCount := o.requestCountInWindow.Sum(size)
	return float64(errorCount / totalCount)
}

//GetRecentErrorPercentage ...
func (o *Stats)GetRecentErrorPercentage() float64 {
	return o.GetErrorPercentage(30)
}