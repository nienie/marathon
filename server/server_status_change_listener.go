package server

//StatusChangeListener server status changed listener
type StatusChangeListener interface {
	//ServerStatusChanged invoked by BaseLoadBalancer when server status has changed
	// (e.g. when marked as down or found dead by ping).
	// parameters represents the servers that had their status changed, never nil
	ServerStatusChanged([]*Server)
}
