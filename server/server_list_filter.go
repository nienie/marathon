package server

//ListFilter This interface allows for filtering the configured or dynamically obtained
//List of candidate servers with desirable characteristics.
type ListFilter interface {
    //GetFilteredListOfServers ...
    GetFilteredListOfServers([]*Server) []*Server
}