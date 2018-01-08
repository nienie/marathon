package loadbalancer

import (
	"context"
	"fmt"

	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/errors"
	"github.com/nienie/marathon/loadbalancer/command"
	"github.com/nienie/marathon/loadbalancer/retry"
	"github.com/nienie/marathon/ratelimit"
	"github.com/nienie/marathon/server"
	"github.com/nienie/marathon/metric"
)

//BaseLoadBalancerClient a default implementation that provide the integration of client with load balancer.
type BaseLoadBalancerClient struct {
	*Context
	client.Client
}

//ExecuteWithLoadBalancer ...
func (c *BaseLoadBalancerClient) ExecuteWithLoadBalancer(ctx context.Context, request client.Request, requestConfig config.ClientConfig) (client.Response, error) {
	if request == nil {
		return nil, errors.NewClientError(errors.General, fmt.Errorf("invalid parameters, request is nil"))
	}
	loadBalancerCommand := c.buildLoadBalancerCommand(request, requestConfig)
	serverOperation := command.ServerOperation(func(server *server.Server) (client.Response, error) {
		serverStats := c.GetServerStats(server)
		if ratelimit.Allow(request.GetURI(), serverStats, requestConfig) == false {
			return nil, errors.NewClientError(errors.ClientThrottled, nil)
		}
		finalURI, err := c.ReconstructURIWithServer(server, request.GetURI())
		if err != nil {
			return nil, err
		}
		requestForServer := request.ReplaceURI(finalURI)
		watch := metric.NewBasicStopWatch()
		watch.Start()
		response, err := c.Client.Execute(ctx, requestForServer, requestConfig)
		watch.Stop()
		metric.RPC(ctx, requestForServer, response, err, watch.GetDuration())
		return response, err
	})
	return loadBalancerCommand.Execute(ctx, serverOperation)
}

func (c *BaseLoadBalancerClient) buildLoadBalancerCommand(request client.Request, requestConfig config.ClientConfig) *Command {
	cmd := NewCommand()
	cmd.WithLoadBalancer(c.LoadBalancer)
	cmd.WithLoadBalancerContext(c.Context)
	cmd.WithServerLocator(request.GetLoadBalancerKey())
	cmd.WithLoadBalancerURI(request.GetURI())
	cmd.WithRetryHandler(c.getRequestSpecificRetryHandler(request, requestConfig))
	return cmd
}

func (c *BaseLoadBalancerClient) getRequestSpecificRetryHandler(request client.Request, requestConfig config.ClientConfig) retry.Handler {
	return retry.NewHTTPClientLoadBalancerRetryHandler(requestConfig)
}
