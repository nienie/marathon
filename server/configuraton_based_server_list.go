package server

import "github.com/nienie/marathon/config"

type ConfigurationBasedServerList struct {
	clientConfig config.ClientConfig
}

func (l *ConfigurationBasedServerList) GetInitialListOfServers() []*Server {
	return l.GetUpdatedListOfServers()
}

func (l *ConfigurationBasedServerList) GetUpdatedListOfServers() []*Server {
	ret, _ := ParseServerListString(l.clientConfig.GetPropertyAsString(config.ListOfServers, config.DefaultListOfServers))
	return ret
}
