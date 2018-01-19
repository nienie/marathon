package server

import (
	"fmt"
)

const (
	//ClusterUnknown ...
	ClusterUnknown = "unknown"
	//DefaultWight ...
	DefaultWight = 10
)

//Server represents a typical server, use Host:Port identifier
type Server struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Scheme      string `json:"scheme"`
	IsAliveFlag bool   `json:"is_alive"`
	TempDown    bool   `json:"-"`
	Cluster     string `json:"cluster"`
	Weight      int    `json:"weight"`
}

//NewServer create a server instance
func NewServer(scheme string, host string, port int) *Server {
	return &Server{
		Scheme:      scheme,
		Host:        host,
		Port:        port,
		IsAliveFlag: true,
		TempDown:    false,
		Cluster:     ClusterUnknown,
		Weight:      DefaultWight,
	}
}

//SetHost ...
func (s *Server) SetHost(host string) *Server {
	s.Host = host
	return s
}

//GetHost ...
func (s *Server) GetHost() string {
	return s.Host
}

//SetPort ...
func (s *Server) SetPort(port int) *Server {
	s.Port = port
	return s
}

//GetPort ...
func (s *Server) GetPort() int {
	return s.Port
}

//GetHostPort ...
func (s *Server) GetHostPort() string {
	if s.Port <= 0 {
		return s.Host
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

//GetScheme ...
func (s *Server) GetScheme() string {
	return s.Scheme
}

//SetAlive ...
func (s *Server) SetAlive(isAliveFlag bool) *Server {
	s.IsAliveFlag = isAliveFlag
	return s
}

//IsAlive ...
func (s *Server) IsAlive() bool {
	return s.IsAliveFlag
}

//Equals ...
func (s *Server) Equals(ss *Server) bool {
	return s.GetHostPort() == ss.GetHostPort() && s.GetScheme() == ss.GetScheme()
}

//GetCluster ...
func (s *Server) GetCluster() string {
	return s.Cluster
}

//SetCluster ...
func (s *Server) SetCluster(cluster string) *Server {
	if len(cluster) == 0 {
		cluster = ClusterUnknown
	}
	s.Cluster = cluster
	return s
}

//SetTempDown ...
func (s *Server) SetTempDown(isDown bool) *Server {
	s.TempDown = isDown
	return s
}

//IsTempDown ...
func (s *Server) IsTempDown() bool {
	return s.TempDown
}

//GetWeight ...
func (s *Server) GetWeight() int {
	return s.Weight
}

//SetWeight ...
func (s *Server) SetWeight(weight int) *Server {
	s.Weight = weight
	return s
}
