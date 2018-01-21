package server

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestNewDefaultServerStatsStats(t *testing.T) {
    ss := NewDefaultServerStats()
    assert.NotNil(t, ss)

    svr := NewServer("http", "127.0.0.1", 8080)
    assert.NotNil(t, svr)

    ss.Initialize(svr)

    ss.IncrementActiveRequestsCount()
    ss.IncrementNumRequests()
    ss.IncrementOpenConnectionsCount()
    ss.AddToFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()
    ss.IncrementSuccessiveConnectionFailureCount()

    for i := 0; i < 100; i++ {
        ss.NoteResponseTime(float64(i))
    }

    assert.Equal(t, int64(1), ss.GetFailureCount())
    assert.Equal(t, int64(1), ss.GetActiveRequestsCount(time.Duration(time.Now().UnixNano())))
    assert.Equal(t, int64(1), ss.GetTotalRequestsCount())
    assert.Equal(t, int64(1), ss.GetOpenConnectionsCount())
    assert.Equal(t, true, ss.IsCircuitBreakerTripped(time.Duration(time.Now().UnixNano())))

    assert.Equal(t, float64(99), ss.GetResponseTimeMax())
    assert.Equal(t, float64(0), ss.GetResponseTimeMin())
    assert.Equal(t, float64(49.5), ss.GetResponseTimeAvg())
    assert.Equal(t, float64(49.5), ss.GetResponseTimeAvgRecent())
    assert.Equal(t, float64(49.5), ss.GetResponseTime50thPercentile())

    assert.Equal(t, int64(1), ss.GetMeasuredRequestsCount())
}