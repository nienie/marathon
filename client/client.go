package client

import (
	"github.com/nienie/marathon/config"
	"context"
)

//Client ...
type Client interface {
	//Execute ...
	Execute(context.Context, Request, config.ClientConfig) (Response, error)
}
