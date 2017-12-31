package server

import (
    "sync/atomic"
    "time"

    "github.com/nienie/marathon/stats"
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
    DefaultConnectionFailureCountThreshold = 3
    //DefaultCircuitTrippedTimeoutFactor ...
    DefaultCircuitTrippedTimeoutFactor = 10
    //DefaultMaxCircuitTrippedTimeout ...
    DefaultMaxCircuitTrippedTimeout = 30 * time.Second
    //DefaultBufferSize ...
    DefaultBufferSize = 60 * 1000 // = 1000 requests / sec for 1 minute
    //DefaultPublishInterval ...
    DefaultPublishInterval = 60 * time.Second // = 1 minute
    //DefaultActiveRequestsCountTimeout ...
    DefaultActiveRequestsCountTimeout = 60 * time.Second
    //DefaultFailureCountSlidingWindowInterval ...
    DefaultFailureCountSlidingWindowInterval = 1 * time.Second
    //DefaultRequestCountsSlidingWindowInterval ...
    DefaultRequestCountsSlidingWindowInterval = 60 * time.Second
)

var (
    percentiles = []float64{
        tenPercentile, twentyFivePercentile, fiftyPercentile, seventyFivePercentile, ninetyPercentile,
        ninetyFivePercentile, ninetyEightPercentile, ninetyNinePercentile, ninetyNinePointFivePercentile,
    }
)

//Stats ...
type Stats struct {
    Server                                 *Server

    ConnectionFailureThreshold             int
    CircuitTrippedTimeoutFactor            int
    MaxCircuitTrippedTimeout               time.Duration
    FailureCountSlidingWindowInterval      time.Duration
    ActiveRequestsCountTimeout             time.Duration
    RequestCountsSlidingWindowInterval     time.Duration

    //for stats
    totalRequests                          int64
    activeRequestsCount                    int64
    openConnectionsCount                   int64
    successiveConnectionFailureCount       int64
    totalCircuitBreakerBlackOutPeriod      int64  //nanoseconds

    //record time
    lastConnectionFailedTimestamp          time.Duration
    lastActiveRequestsCountChangeTimestamp time.Duration
    firstConnectionTimestamp               time.Duration
    lastAccessedTimestamp                  time.Duration

    //stats objects
    responseTimeDist                       *stats.Distribution     //to stats in the overall time
    publisher                              *stats.DataPublisher
    dataDist                               *stats.DataDistribution  // to stats in a recent time-slice
    serverFailureCounts                    *stats.MeasureRate   //server failure counts in a sliding window time
    requestCountInWindow                   *stats.MeasureRate  //request count in a window time

    //stats objects initial parameters
    BufferSize                             int
    PublishInterval                        time.Duration
}

//NewDefaultServerStats ...
func NewDefaultServerStats() *Stats {
    return &Stats{
        ConnectionFailureThreshold:         DefaultConnectionFailureCountThreshold,
        CircuitTrippedTimeoutFactor:        DefaultCircuitTrippedTimeoutFactor,
        MaxCircuitTrippedTimeout:           DefaultMaxCircuitTrippedTimeout,
        FailureCountSlidingWindowInterval:  DefaultFailureCountSlidingWindowInterval,
        ActiveRequestsCountTimeout:         DefaultActiveRequestsCountTimeout,
        RequestCountsSlidingWindowInterval: DefaultRequestCountsSlidingWindowInterval,
        BufferSize:                         DefaultBufferSize,
        PublishInterval:                    DefaultPublishInterval,
        responseTimeDist:                   stats.NewDistribution(),
    }
}

//Initialize ...
func (o *Stats)Initialize(svr *Server) {
    o.serverFailureCounts = stats.NewMeasureRate(o.FailureCountSlidingWindowInterval)
    o.requestCountInWindow = stats.NewMeasureRate(o.RequestCountsSlidingWindowInterval)
    if o.publisher == nil {
        o.dataDist = stats.NewDataDistribution(o.BufferSize, percentiles)
        o.publisher = stats.NewDataPublisher(o.dataDist, o.PublishInterval)
        o.publisher.Start()
    }
    o.Server = svr
}

//Close ...
func (o *Stats)Close() {
    if o.publisher != nil {
        o.publisher.Stop()
    }
}

//AddToFailureCount increment the count of failure for this server
func (o *Stats) AddToFailureCount() {
    o.serverFailureCounts.Increment()
}

//GetFailureCount returns the count of failures in the current window.
func (o *Stats) GetFailureCount() int64 {
    return o.serverFailureCounts.GetCurrentCount()
}

//NoteResponseTime call this method to note the response time after every request.
func (o *Stats) NoteResponseTime(msecs float64) {
    o.dataDist.NoteValue(msecs)
    o.responseTimeDist.NoteValue(msecs)
}

//IncrementNumRequests note the total number of requests.
func (o *Stats) IncrementNumRequests() {
    atomic.AddInt64(&o.totalRequests, 1)
}

//IncrementActiveRequestsCount note the active number of requests.
func (o *Stats) IncrementActiveRequestsCount() {
    atomic.AddInt64(&o.activeRequestsCount, 1)
    o.requestCountInWindow.Increment()
    currentTime := time.Duration(time.Now().UnixNano())
    o.lastActiveRequestsCountChangeTimestamp = currentTime
    o.lastAccessedTimestamp = currentTime
    if o.firstConnectionTimestamp == time.Duration(0) {
        o.firstConnectionTimestamp = currentTime
    }
}

//DecrementActiveRequestsCount ...
func (o *Stats) DecrementActiveRequestsCount() {
    if atomic.AddInt64(&o.activeRequestsCount, -1) < int64(0) {
        atomic.StoreInt64(&o.activeRequestsCount, 0)
    }
    o.lastActiveRequestsCountChangeTimestamp = time.Duration(time.Now().UnixNano())
}

//GetActiveRequestsCount ...
func (o *Stats) GetActiveRequestsCount() int64 {
    count := atomic.LoadInt64(&o.activeRequestsCount)

    currentTime := time.Duration(time.Now().UnixNano())
    if currentTime - o.lastActiveRequestsCountChangeTimestamp > o.ActiveRequestsCountTimeout || count < 0 {
        atomic.StoreInt64(&o.activeRequestsCount, 0)
        return 0
    }

    return count
}

//IncrementOpenConnectionsCount ...
func (o *Stats)IncrementOpenConnectionsCount() {
    atomic.AddInt64(&o.openConnectionsCount, 1)
}

//DecrementOpenConnectionsCount ...
func (o *Stats)DecrementOpenConnectionsCount() {
    if atomic.AddInt64(&o.openConnectionsCount, -1) < 0 {
        atomic.StoreInt64(&o.openConnectionsCount, 0)
    }
}

//GetOpenConnectionsCount ...
func (o *Stats)GetOpenConnectionsCount() int64 {
    return atomic.LoadInt64(&o.openConnectionsCount)
}

//GetMeasuredRequestsCount ...
func (o *Stats)GetMeasuredRequestsCount() int64 {
    return o.requestCountInWindow.GetCount()
}

//GetMonitoredActiveRequestsCount ...
func (o *Stats)GetMonitoredActiveRequestsCount() int64 {
    return atomic.LoadInt64(&o.activeRequestsCount)
}

func (o *Stats)getCircuitBreakerBlackoutPeriod() time.Duration {
    failureCount := atomic.LoadInt64(&o.successiveConnectionFailureCount)
    if failureCount < int64(o.ConnectionFailureThreshold) {
        return time.Duration(0)
    }

    diff := failureCount -int64(o.ConnectionFailureThreshold)
    if diff > 16 {
        diff = 16
    }

    blackOutSeconds := time.Duration(int64(o.CircuitTrippedTimeoutFactor) * diff * 2) * time.Second
    if blackOutSeconds > o.MaxCircuitTrippedTimeout {
        blackOutSeconds = o.MaxCircuitTrippedTimeout
    }

    return blackOutSeconds
}

func (o *Stats)getCircuitBreakerTimeout() time.Duration {
    blackOutPeriod := o.getCircuitBreakerBlackoutPeriod()

    if blackOutPeriod <= 0 {
        return time.Duration(0)
    }

    return o.lastConnectionFailedTimestamp + blackOutPeriod
}

//IsCircuitBreakerTripped ...
func (o *Stats)IsCircuitBreakerTripped() bool {
    currentTime := time.Duration(time.Now().UnixNano())
    circuitBreakerTimeout := o.getCircuitBreakerTimeout()
    if circuitBreakerTimeout <= 0 {
        return false
    }
    return circuitBreakerTimeout > currentTime
}

//IncrementSuccessiveConnectionFailureCount ...
func (o *Stats)IncrementSuccessiveConnectionFailureCount() {
    o.lastConnectionFailedTimestamp = time.Duration(time.Now().UnixNano())
    atomic.AddInt64(&o.successiveConnectionFailureCount, 1)
    atomic.AddInt64(&o.totalCircuitBreakerBlackOutPeriod, int64(o.getCircuitBreakerBlackoutPeriod()))
}

//ClearSuccessiveConnectionFailureCount ...
func (o *Stats)ClearSuccessiveConnectionFailureCount() {
    atomic.StoreInt64(&o.successiveConnectionFailureCount, 0)
}

//GetSuccessiveConnectionCount ...
func (o *Stats)GetSuccessiveConnectionCount() int64 {
    return atomic.LoadInt64(&o.successiveConnectionFailureCount)
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
    return o.dataDist.GetMean()
}

func (o *Stats) getResponseTimePercentile(percent float64) float64 {
    length := len(percentiles)
    var idx int
    for idx = 0; idx < length; idx++ {
        if percentiles[idx] == percent {
            break
        }
    }
    return o.dataDist.GetPercentiles()[idx]
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
    return atomic.LoadInt64(&o.totalRequests)
}
