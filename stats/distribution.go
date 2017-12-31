package stats

import "math"

//Distribution Accumulator of statistics about a distribution of observed values that are produced incrementally.
// Distribution implements DataCollector
type Distribution struct {
    numValues       int64
    sumValues       float64
    sumSquareValues float64
    minValue        float64
    maxValue        float64
}

//NewDistribution ...
func NewDistribution() *Distribution {
    return &Distribution{}
}

//NoteValue ...
func (o *Distribution) NoteValue(val float64) {
    o.numValues++
    o.sumValues += val
    o.sumSquareValues += val * val
    if o.numValues == 1 {
        o.minValue = val
        o.maxValue = val
    } else if val < o.minValue {
        o.minValue = val
    } else if val > o.maxValue {
        o.maxValue = val
    }
}

//Clear ...
func (o *Distribution) Clear() {
    o.numValues = int64(0)
    o.sumValues = float64(0.0)
    o.sumSquareValues = float64(0.0)
    o.minValue = float64(0.0)
    o.maxValue = float64(0.0)
}

//GetNumValues ...
func (o *Distribution) GetNumValues() int64 {
    return o.numValues
}

//GetMean ...
func (o *Distribution) GetMean() float64 {
    if o.numValues <= 1 {
        return o.sumValues
    }

    return o.sumValues / float64(o.numValues)
}

//GetVariance ...
func (o *Distribution) GetVariance() float64 {
    if o.numValues < 2 {
        return float64(0.0)
    } else if o.sumValues == 0.0 {
        return float64(0.0)
    } else {
        mean := o.GetMean()
        return o.sumSquareValues / float64(o.numValues) - mean * mean
    }
}

//GetStdDev ...
func (o *Distribution) GetStdDev() float64 {
    return math.Sqrt(o.GetVariance())
}

//GetMinimum ...
func (o *Distribution) GetMinimum() float64 {
    return o.minValue
}

//GetMaximum ...
func (o *Distribution) GetMaximum() float64 {
    return o.maxValue
}

//Add ...
func (o *Distribution) Add(another *Distribution) {
    o.numValues += another.numValues
    o.sumValues += another.sumValues
    o.sumSquareValues += another.sumSquareValues
    if o.minValue > another.minValue {
        o.minValue = another.minValue
    }
    if o.maxValue < another.maxValue {
        o.maxValue = another.maxValue
    }
}

//Clone clone self
func (o *Distribution) Clone() *Distribution {
    distribution := &Distribution{
        numValues:          o.numValues,
        sumValues:          o.sumValues,
        sumSquareValues:    o.sumSquareValues,
        minValue:           o.minValue,
        maxValue:           o.maxValue,
    }
    return distribution
}