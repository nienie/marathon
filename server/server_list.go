package server

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	//Delimiter ...
	Delimiter = ","
)

//List Interface that defines the methods sed to obtain the List of Servers
type List interface {
	//GetInitialListOfServers ...
	GetInitialListOfServers() []*Server

	//GetUpdatedListOfServers Return updated list of servers. This is called say every 30 secs
	//(configurable) by the Loadbalancer's Ping cycle
	GetUpdatedListOfServers() []*Server
}

//CompareServerList compare serverList1 and serverList2 equal.
//when the length of serverList1 and serverList2 is equal, and elements in the
// serverList1 and serverList2 are the same and in the same order.
func CompareServerList(serverList1, serverList2 []*Server) bool {
	if serverList1 == nil && serverList2 == nil {
		return true
	}

	if serverList1 == nil || serverList2 == nil {
		return false
	}
	len1 := len(serverList1)
	len2 := len(serverList2)
	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if serverList1[i] != serverList2[i] {
			return false
		}
	}

	return true
}

//CloneServerList ...
func CloneServerList(serverList []*Server) []*Server {
	if serverList == nil {
		return nil
	}

	list := make([]*Server, len(serverList))
	copy(list, serverList)
	return list
}

//ParseServerListString convert a string like "http://127.0.0.1:8080@cluster1,http://localhost:8080@cluster2" into slice Server
func ParseServerListString(svrListStr string) ([]*Server, error) {
	if len(svrListStr) == 0 {
		return nil, fmt.Errorf("empty string")
	}

	ret := make([]*Server, 0)
	svrList := strings.Split(svrListStr, Delimiter)
	for _, svr := range svrList {
		if len(svr) == 0 {
			continue
		}

		var (
			scheme  string
			host    string
			port    int
			cluster string
		)

		pos := strings.Index(svr, "://")
		if pos != -1 {
			scheme = svr[:pos]
			svr = svr[pos+3:]
		}

		pos = strings.Index(svr, "@")
		if pos != -1 {
			cluster = svr[pos+1:]
			svr = svr[:pos]
		}

		pos = strings.Index(svr, ":")
		if pos != -1 {
			s := svr[pos+1:]
			p, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return nil, err
			}
			port = int(p)
			svr = svr[:pos]
		}
		host = svr

		server := NewServer(scheme, host, port)
		if len(cluster) != 0 {
			server.SetCluster(cluster)
		}
		ret = append(ret, server)
	}
	return ret, nil
}
