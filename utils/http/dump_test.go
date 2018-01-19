package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestDumpRequestBody ...
func TestDumpRequestBody(t *testing.T) {
	buffer := bytes.NewBufferString(`{"name":"nienie","hobby":"marathon"}`)
	req, err := http.NewRequest(http.MethodPost, "https://www.baidu.com", buffer)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	body, err := DumpRequestBody(req)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	fmt.Println(string(body))
	assert.NotNil(t, req.Body)
	b, err := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	fmt.Println(string(b))
}

//TestDumpResponseBody ...
func TestDumpResponseBody(t *testing.T) {
	buffer := bytes.NewBufferString(`{"name":"nienie","hobby":"marathon"}`)
	resp := &http.Response{
		Body: ioutil.NopCloser(buffer),
	}
	assert.NotNil(t, resp)
	body, err := DumpResponseBody(resp)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	fmt.Println(string(body))
	b, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	fmt.Println(string(b))
}
