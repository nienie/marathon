package metric

import (
	"context"
	"time"

	"github.com/nienie/marathon/client"
)

//Collector ...
type Collector interface {

	//RPC ...
	RPC(context.Context, client.Request, client.Response, error, time.Duration)
}
