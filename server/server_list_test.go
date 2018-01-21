package server

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

//TestParseServerListString ...
func TestParseServerListString(t *testing.T) {
    var (
        svrListStr = "http://127.0.0.1:8080|20@cluster1,https://localhost:80@cluster2"
    )
    servers, err := ParseServerListString(svrListStr)
    assert.Nil(t, err)
    assert.NotNil(t, servers)
    assert.Equal(t, 2, len(servers))
    assert.Equal(t, "http", servers[0].GetScheme())
    assert.Equal(t, 8080, servers[0].GetPort())
    assert.Equal(t, "127.0.0.1:8080", servers[0].GetHostPort())
    assert.Equal(t, 20, servers[0].GetWeight())
    assert.Equal(t, "cluster1", servers[0].GetCluster())
    assert.Equal(t, "https", servers[1].GetScheme())
    assert.Equal(t, 80, servers[1].GetPort())
    assert.Equal(t, "localhost:80", servers[1].GetHostPort())
    assert.Equal(t, "cluster2", servers[1].GetCluster())
    assert.Equal(t, 10, servers[1].GetWeight())
}

//TestCompareServerList ...
func TestCompareServerList(t *testing.T) {
    var (
        svrListStr = "http://127.0.0.1:8080@cluster1,https://localhost:80@cluster2"
    )

    servers, err := ParseServerListString(svrListStr)
    servers1, err1 := ParseServerListString(svrListStr)
    assert.Nil(t, err)
    assert.NotNil(t, servers)
    assert.Equal(t, 2, len(servers))

    assert.Nil(t, err1)
    assert.NotNil(t, servers1)
    assert.Equal(t, 2, len(servers1))

    ss := CloneServerList(servers)
    assert.NotNil(t, ss)
    assert.Equal(t, 2, len(ss))

    isSame := CompareServerList(servers, ss)
    assert.Equal(t, true, isSame)
}