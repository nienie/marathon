package server

//List Interface that defines the methods sed to obtain the List of Servers
type List interface {
    //GetInitialListOfServers ...
    GetInitialListOfServers() []*Server

    //GetUpdatedListOfServers Return updated list of servers. This is called say every 30 secs
    //(configurable) by the Loadbalancer's Ping cycle
    GetUpdatedListOfServers() []*Server
}
