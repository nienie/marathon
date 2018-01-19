package loadbalancer

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	"github.com/nienie/marathon/server"
)

//HashRule ...
type HashRule struct {
	BaseRule
}

//NewHashRule ...
func NewHashRule() Rule {
	return &HashRule{}
}

//Choose ...
func (r *HashRule) Choose(key interface{}) *server.Server {
	return r.ChooseFromLoadBalancer(r.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (r *HashRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	if lb == nil {
		return nil
	}

	reachableServers := lb.GetReachableServers()
	allServers := lb.GetAllServers()

	upCount := len(reachableServers)
	serverCount := len(allServers)

	if upCount == 0 || serverCount == 0 {
		return nil
	}

	h := sha1.New()
	h.Write([]byte(fmt.Sprint(key)))
	bt := h.Sum(nil)
	i, err := binary.ReadUvarint(bytes.NewBuffer(bt))
	if err != nil {
		return nil
	}

	selectedServer := reachableServers[i%uint64(upCount)]
	if selectedServer.IsTempDown() {
		return nil
	}

	return selectedServer
}
