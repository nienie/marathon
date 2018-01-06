package loadbalancer

import (
	"github.com/nienie/marathon/client"
	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/loadbalancer/command"
	"github.com/nienie/marathon/loadbalancer/retry"
	"github.com/nienie/marathon/server"
)

//BaseLoadBalancerClient a default implementation that provide the integration of client with load balancer.
type BaseLoadBalancerClient struct {
	*Context
	client.Client
}

//ExecuteWithLoadBalancer ...
func (c *BaseLoadBalancerClient) ExecuteWithLoadBalancer(request client.Request, requestConfig config.ClientConfig) (client.Response, error) {
	loadBalancerCommand := c.buildLoadBalancerCommand(request, requestConfig)

	serverOperation := command.ServerOperation(func(server *server.Server) (client.Response, error) {
		finalURI, err := c.ReconstructURIWithServer(server, request.GetURI())
		if err != nil {
			return nil, err
		}
		requestForServer := request.ReplaceURI(finalURI)
		return c.Client.Execute(requestForServer, requestConfig)
	})
	return loadBalancerCommand.Execute(serverOperation)
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
