package client

import (
	"net/url"
)

//Response ...
type Response interface {
	//GetPayload ...
	GetPayload() ([]byte, error)
	//HasPayload ...
	HasPayload() bool
	//IsSuccess true if the response is deemed success, for example, 200 response code for http protocal
	IsSuccess() bool
	//GetRequestedURI return the request URI that generated this response
	GetRequestedURI() *url.URL
	//GetHeaders ...
	GetHeaders() map[string][]string
}
