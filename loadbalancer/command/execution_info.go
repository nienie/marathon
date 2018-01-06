package command

import (
	"github.com/nienie/marathon/server"
)

//ExecutionInfo ...
type ExecutionInfo struct {
	Server                       *server.Server
	NumberOfPastAttemptOnServer  int
	NumberOfPastServersAttempted int
}

//ExecutionInfoContext ...
type ExecutionInfoContext struct {
	server             *server.Server
	serverAttemptCount int
	attemptCount       int
}

//NewExecutionInfoContext ...
func NewExecutionInfoContext() *ExecutionInfoContext {
	return &ExecutionInfoContext{}
}

//SetServer ...
func (o *ExecutionInfoContext) SetServer(server *server.Server) {
	o.server = server
	o.serverAttemptCount++
	o.attemptCount = 0
}

//GetServer ...
func (o *ExecutionInfoContext) GetServer() *server.Server {
	return o.server
}

//IncAttemptCount ...
func (o *ExecutionInfoContext) IncAttemptCount() {
	o.attemptCount++
}

//GetAttemptCount ...
func (o *ExecutionInfoContext) GetAttemptCount() int {
	return o.attemptCount
}

//GetServerAttemptCount ...
func (o *ExecutionInfoContext) GetServerAttemptCount() int {
	return o.serverAttemptCount
}

//ToExecutionInfo ...
func (o *ExecutionInfoContext) ToExecutionInfo() *ExecutionInfo {
	return &ExecutionInfo{
		Server: o.server,
		NumberOfPastAttemptOnServer:  o.attemptCount - 1,
		NumberOfPastServersAttempted: o.serverAttemptCount - 1,
	}
}

//ToFinalExecutionInfo ...
func (o *ExecutionInfoContext) ToFinalExecutionInfo() *ExecutionInfo {
	return &ExecutionInfo{
		Server: o.server,
		NumberOfPastAttemptOnServer:  o.attemptCount,
		NumberOfPastServersAttempted: o.serverAttemptCount - 1,
	}
}
