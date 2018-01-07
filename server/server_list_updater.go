package server

import "time"

//UpdateAction server list update action
type UpdateAction interface {
	//DoUpdate an interface for the updateAction that actually executes a server list update
	DoUpdate()
}

//ListUpdater strategy for DynamicServerListLoadBalancer to use for different ways
// of doing dynamic server list updates.
type ListUpdater interface {
	//Start start the serverList updater with the given update action
	//This call should be idempotent.
	Start(UpdateAction)

	//Stop stop the serverList updater. This call should be idempotent
	Stop()

	//GetLastUpdateTime the last update timestamp as a Date string
	GetLastUpdateTime() time.Duration
}
