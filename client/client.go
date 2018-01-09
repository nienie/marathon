package client

import (
	"context"

	"github.com/nienie/marathon/config"
)

//Client ...
type Client interface {
	//Execute ...
	Execute(context.Context, Request, config.ClientConfig) (Response, error)
}
