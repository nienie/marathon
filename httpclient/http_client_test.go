package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nienie/marathon/config"
	"github.com/nienie/marathon/loadbalancer"
	"github.com/nienie/marathon/server"

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

//TestNewHTTPLoadBalancerClient ...
func TestNewHTTPLoadBalancerClient(t *testing.T) {

	prop, err := parsePropertiesFile()
	assert.Nil(t, err)
	clientConfig := config.NewDefaultClientConfig("baidu", prop)
	assert.NotNil(t, clientConfig)

	lb := loadbalancer.NewBaseLoadBalancer(clientConfig, nil, nil, nil)
	assert.NotNil(t, lb)

	serverList := server.NewConfigurationBasedServerList(clientConfig)
	assert.NotNil(t, serverList)
	lb.AddServers(serverList.GetInitialListOfServers())

	httpClient := NewHTTPLoadBalancerClient(clientConfig, lb)
	assert.NotNil(t, httpClient)

	body := bytes.NewBufferString(`{"name":"nienie","hobby":"marathon"}`)
	req, err := NewHTTPRequest(http.MethodPost, "/", body, nil)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	resp, err := httpClient.Do(nil, req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	fmt.Println(string(response))
}
