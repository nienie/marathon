package client

import "github.com/nienie/marathon/config"

//Client ...
type Client interface {
	//Execute ...
	Execute(Request, config.ClientConfig) (Response, error)
}
