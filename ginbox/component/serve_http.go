package component

import (
	"github.com/joyous-x/enceladus/common/xlog"
	"io"
	"net"
	"net/http"
	"strings"
)

// ServeHTTP 转发http请求
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	xlog.Debug("ServeHTTP : Received request %s %s %s", req.Method, req.Host, req.RemoteAddr)
	other := "https://baog.xxxx.com" + req.URL.String()

	outReq, err := http.NewRequest(req.Method, other, req.Body)
	outReq.Header = req.Header

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
		xlog.Debug("ServeHTTP : remote:%v, X-Forwarded-For=%s", other, clientIP)
	}

	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		xlog.Error("ServeHTTP : request %s, error: %v", other, err)
		return
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
	xlog.Debug("ServeHTTP OK: Received request %s %s %s", req.Method, req.Host, req.RemoteAddr)
}
