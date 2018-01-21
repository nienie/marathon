package loadbalancer

import (
	"context"
	"net/url"
	"time"

	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/errors"
	"github.com/nienie/marathon/loadbalancer/command"
	"github.com/nienie/marathon/loadbalancer/retry"
	"github.com/nienie/marathon/metric"
	"github.com/nienie/marathon/server"
)

//Command ...
type Command struct {
	LoadBalancerURI     *url.URL
	LoadBalancerKey     interface{}
	LoadBalancerContext *Context
	LoadBalancer        LoadBalancer
	RetryHandler        retry.Handler
	Server              *server.Server
}

//NewCommand ...
func NewCommand() *Command {
	return &Command{}
}

//WithLoadBalancer ...
func (c *Command) WithLoadBalancer(loadBalancer LoadBalancer) *Command {
	c.LoadBalancer = loadBalancer
	return c
}

//WithLoadBalancerURI ...
func (c *Command) WithLoadBalancerURI(loadBalancerURI *url.URL) *Command {
	c.LoadBalancerURI = loadBalancerURI
	return c
}

//WithRetryHandler ...
func (c *Command) WithRetryHandler(retryHandler retry.Handler) *Command {
	c.RetryHandler = retryHandler
	return c
}

//WithLoadBalancerContext ...
func (c *Command) WithLoadBalancerContext(loadBalancerContext *Context) *Command {
	c.LoadBalancerContext = loadBalancerContext
	return c
}

//WithServerLocator ...
func (c *Command) WithServerLocator(key interface{}) *Command {
	c.LoadBalancerKey = key
	return c
}

//WithServer ...
func (c *Command) WithServer(server *server.Server) *Command {
	c.Server = server
	return c
}

//SelectServer ...
func (c *Command) SelectServer() (*server.Server, error) {
	if c.Server != nil {
		return c.Server, nil
	}
	return c.LoadBalancerContext.GetServerFromLoadBalancer(c.LoadBalancerURI, c.LoadBalancerKey)
}

//Execute ...
func (c *Command) Execute(ctx context.Context, serverOperation command.ServerOperation) (response client.Response, err error) {
	exeCtx := command.NewExecutionInfoContext()
	maxRetrySame := c.RetryHandler.GetMaxRetriesOnSameServer()
	maxRetryNext := c.RetryHandler.GetMaxRetriesOnNextServer()

	server, err := c.SelectServer()
	if err != nil {
		return nil, err
	}
	exeCtx.SetServer(server)

	response, err = c.execute(ctx, exeCtx, server, serverOperation)
	if err == nil {
		return response, err
	}

	//retry on the same server
	if maxRetrySame > 0 {
		for {
			response, err = c.execute(ctx, exeCtx, server, serverOperation)
			if err == nil {
				return response, err
			}

			retryChecker := c.retryPolicy(maxRetrySame, true)
			if !retryChecker(exeCtx.GetAttemptCount(), err) {
				break
			}
		}
	}

	if maxRetrySame > 0 && maxRetryNext == 0 && exeCtx.GetAttemptCount() == (maxRetrySame+1) {
		return nil, errors.NewClientError(errors.NumberOfRetriesExceeded, err)
	}

	//retry on different serverss
	if maxRetryNext > 0 && c.Server == nil {
		for {
			server, err = c.SelectServer()
			if err != nil {
				return nil, err
			}
			exeCtx.SetServer(server)

			response, err = c.execute(ctx, exeCtx, server, serverOperation)
			if err == nil {
				return response, err
			}

			retryChecker := c.retryPolicy(maxRetryNext, false)
			if !retryChecker(exeCtx.GetServerAttemptCount(), err) {
				break
			}
		}
	}

	if maxRetryNext > 0 && exeCtx.GetServerAttemptCount() == (maxRetryNext+1) {
		return nil, errors.NewClientError(errors.NumberOfRetriesNextServerExceeded, err)
	}

	return response, err
}

func (c *Command) execute(ctx context.Context, exeCtx *command.ExecutionInfoContext, server *server.Server, operation command.ServerOperation) (client.Response, error) {
	exeCtx.IncAttemptCount()
	stats := c.LoadBalancerContext.GetServerStats(server)
	c.LoadBalancerContext.NoteOpenConnection(stats)
	stopWatch := metric.NewBasicStopWatch()
	stopWatch.Start()
	response, err := operation(server)
	stopWatch.Stop()
	c.recordStats(ctx, stats, response, err, stopWatch.GetDuration())
	return response, err
}

func (c *Command) recordStats(ctx context.Context, stats *server.Stats, response client.Response, err error, responseTime time.Duration) {
	c.LoadBalancerContext.NoteRequestCompletion(ctx, stats, response, err, int64(responseTime/time.Millisecond), c.RetryHandler)
}

func (c *Command) retryPolicy(maxRetries int, same bool) command.RetryChecker {
	retryCheck := func(tryCount int, err error) bool {
		switch err.(type) {
		case errors.ClientError:
			errorType := err.(errors.ClientError).GetErrType()
			if errorType == errors.AbortExecutionException {
				return false
			}
		default:
		}

		if tryCount > maxRetries {
			return false
		}

		return c.RetryHandler.IsRetriableException(err, same)
	}
	return command.RetryChecker(retryCheck)
}
