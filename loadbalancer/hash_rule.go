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
func (o *HashRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *HashRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	if lb == nil {
		return nil
	}

	upList := lb.GetReachableServers()
	upCount := len(upList)
	if upCount == 0 {
		return nil
	}

	h := sha1.New()
	h.Write([]byte(fmt.Sprint(key)))
	bt := h.Sum(nil)
	i, err := binary.ReadUvarint(bytes.NewBuffer(bt))
	if err != nil {
		return nil
	}

	return upList[i%uint64(upCount)]
}
