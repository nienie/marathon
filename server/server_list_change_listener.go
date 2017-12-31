package server

//ListChangeListener when server list is changed
type ListChangeListener interface {
    //ServerListChanged invoke by BaseLoadBalancer when server list is changed
    ServerListChanged(oldList []*Server, newList []*Server)
}