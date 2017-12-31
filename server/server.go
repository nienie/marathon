package server

import (
    "fmt"
    "strconv"
    "strings"
)

const (
    //IDFormat server identifier format
    IDFormat = "%s:%d"
    //ClusterUnknown ...
    ClusterUnknown = "unknown"
)

//Server represents a typical server, use Host:Port identifier
type Server struct {
    id           string
    host         string
    port         int
    scheme       string
    isAliveFlag  bool
    readyToServe bool
    cluster      string
}

//NewServer create a server instance
func NewServer(scheme string, host string, port int) *Server {
    s := &Server{
        scheme:      scheme,
        host:        host,
        port:        port,
        id:          fmt.Sprintf(IDFormat, host, port),
        isAliveFlag: false,
        cluster:     ClusterUnknown,
    }
    return s
}

//SetID ...
func (s *Server) SetID(id string) {
    scheme, host, port, err := parseID(id)
    if err != nil {
        return
    }

    s.scheme = scheme
    s.host = host
    s.port = port

    s.id = fmt.Sprintf(IDFormat, host, port)
    return
}

func parseID(id string) (scheme string, host string, port int, err error) {
    if len(id) == 0 {
        return
    }

    if strings.HasPrefix(id, "http://") {
        id = strings.TrimLeft(id, "http://")
        scheme = "http"
    } else if strings.HasPrefix(id, "https://") {
        id = strings.TrimLeft(id, "https://")
        scheme = "https"
    }

    if strings.Contains(id, "/") {
        slashIdx := strings.Index(id, "/")
        id = id[:slashIdx]
    }

    colonIdx := strings.Index(id, ":")
    if colonIdx == -1 {
        host = id
        port = 80
        return
    }

    host = id[:colonIdx]
    p, err := strconv.ParseInt(id[colonIdx+1:], 10, 32)
    if err != nil {
        return
    }
    port = int(p)

    return
}

//GetID ...
func (s *Server) GetID() string {
    return s.id
}

//SetHost ...
func (s *Server) SetHost(host string) {
    s.host = host
    s.id = fmt.Sprintf(IDFormat, s.host, s.port)
}

//GetHost ...
func (s *Server) GetHost() string {
    return s.host
}

//SetPort ...
func (s *Server) SetPort(port int) {
    s.port = port
    s.id = fmt.Sprintf(IDFormat, s.host, s.port)
}

//GetPort ...
func (s *Server) GetPort() int {
    return s.port
}

//GetHostPort ...
func (s *Server) GetHostPort() string {
    return fmt.Sprintf(IDFormat, s.host, s.port)
}

//GetScheme ...
func (s *Server) GetScheme() string {
    return s.scheme
}

//SetAlive ...
func (s *Server) SetAlive(isAliveFlag bool) {
    s.isAliveFlag = isAliveFlag
}

//IsAlive ...
func (s *Server) IsAlive() bool {
    return s.isAliveFlag
}

//SetReadyToServe ...
func (s *Server) SetReadyToServe(readyToServe bool) {
    s.readyToServe = readyToServe
}

//IsReadyToServe ...
func (s *Server) IsReadyToServe() bool {
    return s.readyToServe
}

//Equals ...
func (s *Server) Equals(ss *Server) bool {
    return s.GetID() == ss.GetID()
}

//GetCluster ...
func (s *Server)GetCluster() string{
    return s.cluster
}

//SetCluster ...
func (s *Server)SetCluster(cluster string){
    if len(cluster) > 0 {
        s.cluster = cluster
    }
}