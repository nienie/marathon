package loadbalancer

//ClusterSnapshot ...
type ClusterSnapshot struct {
	InstanceCount       int
	LoadPerServer       float64
	CircuitTrippedCount int
	ActiveRequestsCount int64
}

//NewDefaultClusterSnapshot ...
func NewDefaultClusterSnapshot() *ClusterSnapshot {
	return &ClusterSnapshot{}
}

//NewClusterSnapshot ...
func NewClusterSnapshot(instanceCount int, loadPerServer float64, circuitTrippedCount int, activeRequestsCount int64) *ClusterSnapshot {
	return &ClusterSnapshot{
		InstanceCount:       instanceCount,
		LoadPerServer:       loadPerServer,
		CircuitTrippedCount: circuitTrippedCount,
		ActiveRequestsCount: activeRequestsCount,
	}
}
