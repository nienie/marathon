package ping

import (
	"io/ioutil"
	"net/http"

	"github.com/nienie/marathon/server"
)

//URLPing Ping implementation if you want to do a "health check" kind of ping.
//This will be a real ping. As in a real http/s call is made to this url.
type URLPing struct {
	IsSecure         bool
	PingAppendString string
	ExpectedContent  string
}

//NewURLPing ...
func NewURLPing(isSecure bool, pingAppendString, expectedContent string) Ping {
	return &URLPing{
		IsSecure:         isSecure,
		PingAppendString: pingAppendString,
		ExpectedContent:  expectedContent,
	}
}

//IsAlive ...
func (p *URLPing) IsAlive(svr *server.Server) bool {
	urlStr := ""
	if p.IsSecure {
		urlStr = urlStr + "https://"
	} else {
		urlStr = urlStr + "http://"
	}
	urlStr = urlStr + svr.GetHostPort()
	urlStr = urlStr + p.PingAppendString
	resp, err := http.Get(urlStr)
	if err != nil {
		return false
	}

	if len(p.ExpectedContent) == 0 {
		return resp.StatusCode == http.StatusOK
	}

	defer resp.Body.Close()
	responseContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	if p.ExpectedContent == string(responseContent) {
		return true
	}
	return false
}
