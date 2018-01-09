package errors

//ErrorType ...
type ErrorType int

const (
	//OK ...
	OK ErrorType = iota
	//General ...
	General
	//Configuration ...
	Configuration
	//NumberOfRetriesExceeded ...
	NumberOfRetriesExceeded
	//NumberOfRetriesNextServerExceeded ...
	NumberOfRetriesNextServerExceeded
	//SocketTimeoutException signals that a timeout occurred on a socket read or accept.
	SocketTimeoutException
	//ReadTimeoutException signals that a timeout occurred on a socket read.
	ReadTimeoutException
	//SocketException thrown to indicate that there is an error creating or accessing a socket.
	SocketException
	//UnknownHostException ...
	UnknownHostException
	//ConnectException signals that an error occurred while attempting to connection a socket to a remote address and port.
	ConnectException
	//ClientThrottled ...
	ClientThrottled
	//ServerThrottled ...
	ServerThrottled
	//NoRouteToHostException ...
	NoRouteToHostException
	//CacheMissing ...
	CacheMissing
	//AbortExecutionException ...
	AbortExecutionException
)

var errorTypeNameMap = map[ErrorType]string{
	OK:                                "OK",
	General:                           "General",
	Configuration:                     "Configuration",
	NumberOfRetriesExceeded:           "NumberOfRetriesExceeded",
	NumberOfRetriesNextServerExceeded: "NumberOfRetriesNextServerExceeded",
	SocketTimeoutException:            "SocketTimeoutException",
	ReadTimeoutException:              "ReadTimeoutException",
	SocketException:                   "SocketException",
	UnknownHostException:              "UnknownHostException",
	ConnectException:                  "ConnectException",
	ClientThrottled:                   "ClientThrottled",
	ServerThrottled:                   "ServerThrottled",
	NoRouteToHostException:            "NoRouteToHostException",
	CacheMissing:                      "CacheMissing",
	AbortExecutionException:           "AbortExecutionException",
}

//GetName ...
func (e ErrorType) GetName() string {
	if name, ok := errorTypeNameMap[e]; ok {
		return name
	}
	return "UnkownErrorType"
}
