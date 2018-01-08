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

//RegisterCollector ...
func RegisterCollector(c Collector) {
	metricCollectors = append(metricCollectors, c)
}

//RPC ...
func RPC(ctx context.Context, request client.Request, response client.Response, err error, costTime time.Duration) {
	for _, c := range metricCollectors {
		c.RPC(ctx, request, response, err, costTime)
	}
}
