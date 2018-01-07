package loadbalancer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/errors"
	"github.com/nienie/marathon/loadbalancer/retry"
	"github.com/nienie/marathon/server"
)

//Context A class contains APIs intended to be used be load balancing client which is subclass of this class.
type Context struct {
	ClientName   string
	LoadBalancer LoadBalancer
	RetryHandler retry.Handler
}

//NewLoadBalancerContext ...
func NewLoadBalancerContext(clientConfig config.ClientConfig, lb LoadBalancer) *Context {
	ctx := &Context{
		ClientName:   clientConfig.GetClientName(),
		LoadBalancer: lb,
		RetryHandler: retry.NewLoadBalancerRetryHandler(clientConfig),
	}

	return ctx
}

func (o *Context) recordStats(stats *server.Stats, responseTime int64) {
	if stats == nil {
		return
	}
	stats.DecrementActiveRequestsCount()
	stats.IncrementNumRequests()
	stats.NoteResponseTime(float64(responseTime))
}

//NoteRequestCompletion This is called after a response is received or an exception is thrown from the client to update related stats.
func (o *Context) NoteRequestCompletion(stats *server.Stats, response interface{},
		err error, responseTime int64, errorHandler retry.Handler) {
	if stats == nil {
		return
	}
	o.recordStats(stats, responseTime)
	callErrorHandler := errorHandler
	if callErrorHandler == nil {
		callErrorHandler = o.RetryHandler
	}

	if err != nil {
		if errorHandler.IsCircuitTrippingException(err) {
			stats.IncrementSuccessiveConnectionFailureCount()
			stats.AddToFailureCount()
			if stats.IsCircuitBreakerTripped(time.Duration(time.Now().UnixNano())) {
				o.LoadBalancer.MarkServerTempDown(stats.Server)
			}
			return
		}
		stats.ClearSuccessiveConnectionFailureCount()
		o.LoadBalancer.MarkServerReady(stats.Server)
		return
	}
	stats.ClearSuccessiveConnectionFailureCount()
	o.LoadBalancer.MarkServerReady(stats.Server)
	return
}

//NoteError This is called after an error is thrown from the client to update related stats.
func (o *Context) NoteError(stats *server.Stats, request client.Request, err error, responseTime int64) {
	if stats != nil {
		return
	}
	o.recordStats(stats, responseTime)
	errorHandler := o.RetryHandler
	if err != nil {
		if errorHandler.IsCircuitTrippingException(err) {
			stats.IncrementSuccessiveConnectionFailureCount()
			stats.AddToFailureCount()
			if stats.IsCircuitBreakerTripped(time.Duration(time.Now().UnixNano())) {
				o.LoadBalancer.MarkServerTempDown(stats.Server)
			}
			return
		}
		stats.ClearSuccessiveConnectionFailureCount()
		o.LoadBalancer.MarkServerReady(stats.Server)
		return
	}

	return
}

//NoteOpenConnection This is usually called just before client execute a request.
func (o *Context) NoteOpenConnection(stats *server.Stats) {
	if stats == nil {
		return
	}
	stats.IncrementActiveRequestsCount()
}

func (o *Context) deriveSchemeAndPortFromPartialURI(uri *url.URL) (port int, scheme string) {
	isSecure := false
	scheme = uri.Scheme
	if len(scheme) != 0 {
		isSecure = strings.Contains(scheme, "https")
	}
	hostPost := strings.Split(uri.Host, ":")
	if len(hostPost) == 2 {
		p, err := strconv.ParseInt(hostPost[1], 10, 32)
		if err != nil {
			p = 0
		}
		port = int(p)
	}
	if port <= 0 && !isSecure {
		port = 80
	} else if port <= 0 && isSecure {
		port = 443
	}
	if len(scheme) == 0 {
		if isSecure {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return
}

//GetServerFromLoadBalancer compute the final URI from a partial URI in the request.
func (o *Context) GetServerFromLoadBalancer(original *url.URL, loadBalancerKey interface{}) (*server.Server, error) {
	var (
		scheme string
		host   string
		port   int
	)
	if original != nil {
		hostPort := strings.Split(original.Host, ":")
		host = hostPort[0]
		port, scheme = o.deriveSchemeAndPortFromPartialURI(original)
	}
	lb := o.LoadBalancer
	//Various Supported Cases
	// The loadbalancer to use and the instances it has is based on how it was registerd
	// In each of these cases, the client might come in use Full Url or Partial Url.
	if len(host) == 0 {
		//Partial URI or no URI case
		//well we have to just get the right instance from lb - or we fall back
		if lb != nil {
			svc := lb.ChooseServer(loadBalancerKey)
			if svc == nil {
				return nil, errors.NewClientError(errors.General, fmt.Errorf("Load balancer does not have available server for client: %s", o.ClientName))
			}
			host = svc.GetHost()
			if len(host) == 0 {
				return nil, errors.NewClientError(errors.General, fmt.Errorf("Invalid Server for : %+v", svc))
			}
			return svc, nil
		}

		return nil, errors.NewClientError(errors.General, fmt.Errorf("Request contains no Host to talk to"))
	}
	return server.NewServer(scheme, host, port), nil
}

//ReconstructURIWithServer ...
func (o *Context) ReconstructURIWithServer(svr *server.Server, original *url.URL) (*url.URL, error) {
	host := svr.GetHost()
	port := svr.GetPort()
	scheme := svr.GetScheme()

	if original.Scheme == scheme && original.Host == svr.GetHostPort() {
		return original, nil
	}

	if len(scheme) == 0 {
		scheme = original.Scheme
	}

	if len(scheme) == 0 {
		_, scheme = o.deriveSchemeAndPortFromPartialURI(original)
	}

	newURL := scheme + "://"

	if original.User != nil {
		newURL = newURL + original.User.Username() + "@"
	}

	newURL = newURL + host

	if port >= 0 {
		newURL = newURL + ":" + fmt.Sprint(port)
	}

	newURL = newURL + original.Path

	if len(original.RawQuery) > 0 {
		newURL = newURL + "?" + original.RawQuery
	}

	if len(original.Fragment) > 0 {
		newURL = newURL + "#" + original.Fragment
	}

	return url.Parse(newURL)
}

//GetServerStats ...
func (o *Context) GetServerStats(svr *server.Server) *server.Stats {
	lb := o.LoadBalancer
	if lb != nil {
		lbStats := lb.GetLoadBalancerStats()
		serverStats := lbStats.GetSingleServerStats(svr)
		return serverStats
	}
	return nil
}
