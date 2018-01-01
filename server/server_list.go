package server

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