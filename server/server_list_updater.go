package server

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
	Start(UpdateAction UpdateAction)

	//Stop stop the serverList updater. This call should be idempotent
	Stop()

	//GetLastUpdate the last update timestamp as a Date string
	GetLastUpdate() string

	//GetDurationSinceLastUpdateMs the number of ms that has elapsed since last update
	GetDurationSinceLastUpdateMs() int64

	//GetNumberMissedCycles the number of update cycles missed, if valid
	GetNumberMissedCycles() int

	//GetCoreThreads the number of threads used, if vaid
	GetCoreThreads() int
}
