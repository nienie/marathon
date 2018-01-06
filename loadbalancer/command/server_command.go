package command

import (
	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/server"
)

//ServerOperation ...
type ServerOperation func(server *server.Server) (client.Response, error)

//RetryChecker ...
type RetryChecker func(tryCount int, err error) bool
