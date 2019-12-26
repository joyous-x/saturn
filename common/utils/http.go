package utils

import (
	"bytes"
	"encoding/json"
	"github.com/joyous-x/saturn/common/xlog"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpPostWwwForm(urls string, data map[string]string) ([]byte, error) {
	vals := url.Values{}
	for k, v := range data {
		vals.Set(k, v)
	}
	body := strings.NewReader(vals.Encode())
	return httpPostBase(urls, body, "application/x-www-form-urlencoded", 5)
}

func HttpPostJson(url string, data interface{}) ([]byte, error) {
	vals, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(vals)
	return httpPostBase(url, body, "application/json", 5)
}

func httpPostBase(url string, body io.Reader, content_type string, conn_timeout_sec int) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(conn_timeout_sec) * time.Second,
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", content_type)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	return d, err
}

// HttpProxy 转发http请求
//   arg: targetHostSchemed, examples "https://baog.xxxx.com"
func HttpProxy(rw http.ResponseWriter, req *http.Request, targetHostSchemed string) {
	xlog.Debug("ServeHTTP : Received request %s %s %s", req.Method, req.Host, req.RemoteAddr)
	other := targetHostSchemed + req.URL.String()

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
