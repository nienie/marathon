package stats

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

//TestNewRollingCounter ...
func TestNewRollingCounter(t *testing.T) {
    counter := NewRollingCounter(1)
    assert.NotNil(t, counter)

    counter.Inc(2)
    assert.Equal(t, int64(2), counter.Count())

    counter.Dec(1)
    assert.Equal(t, int64(1), counter.Count())

    counter.Clear()
    assert.Equal(t, int64(0), counter.Count())

    counter.Inc(10)

    time.Sleep(2 * time.Second)
    assert.Equal(t, int64(0), counter.Count())
}