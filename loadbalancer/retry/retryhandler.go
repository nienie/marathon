package retry

//Handler ...
type Handler interface {

	//IsRetriableException test if an exception is retriable for the load balancer
	IsRetriableException(err error, sameServer bool) bool

	//IsCircuitTrippingException test if an exception should be treated as circuit. For example,
	//a ConnectionException is a circuit failure. This is used to determine whether successive exceptions of
	//such should trip the circuit breaker to a particular host by the load balancer. If false but a server response
	// is absent, load balancer will also cose the circuit upon getting such exception.
	IsCircuitTrippingException(err error) bool

	//GetMaxRetriesOnSameServer Number of maximal retries to be done on one server
	GetMaxRetriesOnSameServer() int

	//GetMaxRetriesOnNextServer Number of maximal different servers to retry
	GetMaxRetriesOnNextServer() int
}
