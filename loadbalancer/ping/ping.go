package ping

import (
	"github.com/nienie/marathon/server"
)

//Ping Interface that defines how we "ping" a server to check if its alive
type Ping interface {
	//IsAlive Checks whether the given Server is "alive"
	IsAlive(*server.Server) bool
}

//NoOpPing ...
type NoOpPing struct{}

//IsAlive ...
func (o *NoOpPing) IsAlive(server *server.Server) bool {
	return true
}
