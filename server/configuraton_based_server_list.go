package server

import "github.com/nienie/marathon/config"

//ConfigurationBasedServerList ...
type ConfigurationBasedServerList struct {
	clientConfig config.ClientConfig
}

//NewConfigurationBasedServerList ...
func NewConfigurationBasedServerList(clientConfig config.ClientConfig) *ConfigurationBasedServerList {
	return &ConfigurationBasedServerList{
		clientConfig: clientConfig,
	}
}

//GetInitialListOfServers ...
func (l *ConfigurationBasedServerList) GetInitialListOfServers() []*Server {
	return l.GetUpdatedListOfServers()
}

//GetUpdatedListOfServers ...
func (l *ConfigurationBasedServerList) GetUpdatedListOfServers() []*Server {
	ret, _ := ParseServerListString(l.clientConfig.GetPropertyAsString(config.ListOfServers, config.DefaultListOfServers))
	return ret
}
