package stats

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

//TestNewRollingSample ...
func TestNewRollingSample(t *testing.T) {
    s := NewRollingSample(1)

    assert.NotNil(t, s)

    for i := int64(0); i < int64(100); i++ {
        s.UpdateValue(i)
    }

    assert.Equal(t, int64(99), s.Max())
    assert.Equal(t, int64(0), s.Min())
    assert.Equal(t, int64(4950), s.Sum())
    assert.Equal(t, float64(49.5), s.Mean())
    assert.Equal(t, float64(49.5), s.Percentile(0.5))
}