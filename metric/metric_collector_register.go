package metric

import (
	"context"
	"time"

	"github.com/nienie/marathon/client"
)

var (
	metricCollectors []Collector
)

func init() {
	metricCollectors = make([]Collector, 0)
}

//RegisterCollectors ...
func RegisterCollectors(cs ...Collector) {
	metricCollectors = append(metricCollectors, cs...)
}

//RPC ...
func RPC(ctx context.Context, request client.Request, response client.Response, err error, costTime time.Duration) {
	for _, c := range metricCollectors {
		if c != nil {
			c.RPC(ctx, request, response, err, costTime)
		}
	}
}
