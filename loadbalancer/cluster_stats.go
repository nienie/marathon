package loadbalancer

//ClusterStats ...
type ClusterStats struct {
	ClusterName       string
	loadBalancerStats *Stats
}

//NewClusterStats ...
func NewClusterStats(cluster string, loadBalanerStats *Stats) *ClusterStats {
	return &ClusterStats{
		ClusterName:       cluster,
		loadBalancerStats: loadBalanerStats,
	}
}

//GetInstanceCount ...
func (o *ClusterStats) GetInstanceCount() int {
	return o.loadBalancerStats.GetInstanceCount(o.ClusterName)
}

//GetCircuitBreakerTrippedCount ...
func (o *ClusterStats) GetCircuitBreakerTrippedCount() int {
	return o.loadBalancerStats.GetCircuitBreakerTrippedCount(o.ClusterName)
}

//GetActiveRequestsPerServer ...
func (o *ClusterStats) GetActiveRequestsPerServer() float64 {
	return o.loadBalancerStats.GetActiveRequestsPerServer(o.ClusterName)
}

//GetMeasuredClusterHits ...
func (o *ClusterStats) GetMeasuredClusterHits() int64 {
	return o.loadBalancerStats.GetMeasuredClusterHits(o.ClusterName)
}
