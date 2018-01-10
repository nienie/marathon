package config

import (
    "fmt"
    "path/filepath"
    "testing"
    "runtime"
    "time"

    "github.com/magiconair/properties"
    "github.com/stretchr/testify/assert"
)

func parsePropertiesFile() (*properties.Properties, error) {
    _, fn, _, _ := runtime.Caller(-1)
    file := filepath.Join(filepath.Dir(fn), "/test_data/test.properties")
    fmt.Println(file)
    prop, err := properties.LoadFile(file, properties.UTF8)
    return prop, err
}

//TestNewDefaultClientConfig ...
func TestNewDefaultClientConfig(t *testing.T) {
    prop, err := parsePropertiesFile()
    if err != nil {
        t.Fatal(err)
    }
    clientConfig := NewDefaultClientConfig("test", prop)
    if clientConfig == nil {
        t.Fatal(fmt.Errorf("NewDefaultClientConfig Failed"))
    }
    b := clientConfig.GetPropertyAsBool(ConcurrencyRateLimitSwitch, false)
    assert.Equal(t, true, b)
    s := clientConfig.GetPropertyAsString(ListOfServers, "")
    assert.Equal(t, "http://127.0.0.1:8000@cluster1,http://localhost:8000@cluster2", s)
    i := clientConfig.GetPropertyAsInteger(MaxAutoRetries, 0)
    assert.Equal(t, 2, i)
    d := clientConfig.GetPropertyAsDuration(ListOfServersPollingInterval, 10 * time.Second)
    assert.Equal(t, 20 * time.Second, d)
    o := clientConfig.GetPropertyAsInteger(MaxAutoRetriesNextServer, 3)
    assert.Equal(t, 1, o)
    clientConfig.SetProperty(MaxTotalConnections, 20)
    c := clientConfig.GetPropertyAsInteger(MaxTotalConnections, 30)
    assert.Equal(t, 20, c)
    name := clientConfig.GetClientName()
    assert.Equal(t, "test", name)
}