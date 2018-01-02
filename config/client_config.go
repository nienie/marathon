package config

import (
    "time"

    "github.com/magiconair/properties"
)

//ClientConfig ...
type ClientConfig interface {
    //GetClientName ...
    GetClientName() string

    //LoadProperties load the properties for a given client and/or load balancer.
    LoadProperties(clientName string)

    //LoadDefaultValues ...
    LoadDefaultValues()

    //GetPropertyAsInteger ...
    GetPropertyAsInteger(configKey string, defaultValue int) int

    //GetPropertyAsString ...
    GetPropertyAsString(configKey string, defaultValue string) string

    //GetPropertyAsBool ...
    GetPropertyAsBool(configKey string, defaultValue bool) bool

    //GetPropertyAsDuration ...
    GetPropertyAsDuration(configKey string, defaultValue time.Duration) time.Duration

    //Set ...
    SetProperty(configKey string, value interface{})
}

//DefaultClientConfig ...
type DefaultClientConfig struct {
    InternalProperties     *properties.Properties
    ExternalProperties     *properties.Properties
    clientName             string
    namespace              string
}

//NewDefaultClientConfig ...
func NewDefaultClientConfig(props *properties.Properties) ClientConfig {
    return &DefaultClientConfig{
        InternalProperties: properties.NewProperties(),
        ExternalProperties: props,
    }
}

//GetNamespace ...
func (c *DefaultClientConfig)GetNamespace() string {
    return c.namespace
}

//SetClientName ...
func (c *DefaultClientConfig)SetClientName(clientName string) {
    c.clientName = clientName
}

//GetClientName ...
func (c *DefaultClientConfig)GetClientName() string {
    return c.clientName
}

func (c *DefaultClientConfig)setPropertyInternal(propName string, value interface{}) {
    c.InternalProperties.SetValue(propName, value)
    return
}

func (c *DefaultClientConfig)putDefaultIntegerProperty(propName string, defaultValue int) {
    value := c.ExternalProperties.GetInt(propName, defaultValue)
    c.setPropertyInternal(propName, value)
}

func (c *DefaultClientConfig)putDefaultBoolProperty(propName string, defaultValue bool) {
    value := c.ExternalProperties.GetBool(propName, defaultValue)
    c.setPropertyInternal(propName, value)
}

func (c *DefaultClientConfig)putDefaultStringProperty(propName string, defaultValue string) {
    value := c.ExternalProperties.GetString(propName, defaultValue)
    c.setPropertyInternal(propName, value)
}

func (c *DefaultClientConfig)putDefaultFloat64Property(propName string, defaultValue float64) {
    value := c.ExternalProperties.GetFloat64(propName, defaultValue)
    c.setPropertyInternal(propName, value)
}

func (c *DefaultClientConfig)putDefaultDurationProperty(propName string, defaultValue time.Duration) {
    value := c.ExternalProperties.GetDuration(propName, defaultValue)
    c.setPropertyInternal(propName, value)
}

//LoadDefaultValues ...
func (c *DefaultClientConfig)LoadDefaultValues() {
    c.putDefaultBoolProperty(EnableConnectionPool, DefaultEnableConnectionPool)
    c.putDefaultIntegerProperty(MaxConnectionsPerHost, DefaultMaxConnectionsPerHost)
    c.putDefaultIntegerProperty(MaxTotalConnections, DefaultMaxTotalConnections)
    c.putDefaultDurationProperty(ConnectTimeout, DefaultConnectTimeout)
    c.putDefaultDurationProperty(ReadTimeout, DefaultReadTimeout)
    c.putDefaultIntegerProperty(MaxAutoRetries, DefaultMaxAutoRetries)
    c.putDefaultIntegerProperty(MaxAutoRetriesNextServer, DefaultMaxAutoRetriesNextServer)
    c.putDefaultBoolProperty(OKToRetryOnAllOperations, DefaultOKToRetryOnAllOperations)
    c.putDefaultIntegerProperty(Port, DefaultPort)
    c.putDefaultStringProperty(ListOfServers, DefaultListOfServers)
    c.putDefaultIntegerProperty(ConnectionFailureThreshold, DefaultConnectionFailureThreshold)
    c.putDefaultIntegerProperty(CircuitTrippedTimeoutFactor, DefaultCircuitTrippedTimeoutFactor)
    c.putDefaultDurationProperty(CircuitTripMaxTimeout, DefaultCircuitTripMaxTimeout)
    c.putDefaultDurationProperty(PingInterval, DefaultPingInterval)
}

//LoadProperties ...
func (c *DefaultClientConfig)LoadProperties(clientName string) {
    c.SetClientName(clientName)
    c.LoadDefaultValues()
    props := c.ExternalProperties.FilterStripPrefix(clientName + ".")
    for _, key := range props.Keys() {
        value := props.GetString(key, "")
        c.setPropertyInternal(key, value)
    }
}

//GetPropertyAsInteger ...
func (c *DefaultClientConfig)GetPropertyAsInteger(configKey string, defaultValue int) int {
    return c.InternalProperties.GetInt(configKey, defaultValue)
}

//GetPropertyAsString ...
func (c *DefaultClientConfig)GetPropertyAsString(configKey string, defaultValue string) string {
    return c.InternalProperties.GetString(configKey, defaultValue)
}

//GetPropertyAsBool ...
func (c *DefaultClientConfig)GetPropertyAsBool(configKey string, defaultValue bool) bool {
    return c.InternalProperties.GetBool(configKey, defaultValue)
}

//GetPropertyAsDuration ...
func (c *DefaultClientConfig)GetPropertyAsDuration(configKey string, defaultValue time.Duration) time.Duration {
    return c.InternalProperties.GetDuration(configKey, defaultValue)
}

//SetProperty ...
func (c *DefaultClientConfig)SetProperty(configKey string, value interface{}) {
    c.InternalProperties.SetValue(configKey, value)
}