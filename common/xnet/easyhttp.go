package xnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/joyous-x/saturn/common/xlog"
)

// HTTPOptions options for http request
type HTTPOptions struct {
	Host        string
	UserAgent   string
	ConnTimeout time.Duration
}

// EasyHTTP easy http request
type EasyHTTP struct {
	options HTTPOptions
}

// NewEasyHTTP construct a new EasyHTTP
func NewEasyHTTP() *EasyHTTP {
	return &EasyHTTP{}
}

// Options return a new EasyHTTP with options in  input
func (s *EasyHTTP) Options(options *HTTPOptions) *EasyHTTP {
	tmp := &EasyHTTP{}
	tmp.SetOptions(options)
	return tmp
}

// SetOptions set options for http request
func (s *EasyHTTP) SetOptions(options *HTTPOptions) {
	if nil != options {
		s.options = *options
	} else {
		s.options = HTTPOptions{
			Host:        "",
			UserAgent:   "",
			ConnTimeout: 5 * time.Second,
		}
	}
}

// Get ...
func (s *EasyHTTP) Get(urlPath string) ([]byte, error) {
	return s.httpBase(urlPath, "GET", nil, "", s.options)
}

// PostWwwForm post x-www-form-urlencoded
func (s *EasyHTTP) PostWwwForm(urlPath string, data map[string]string) ([]byte, error) {
	vals := url.Values{}
	for k, v := range data {
		vals.Set(k, v)
	}
	body := strings.NewReader(vals.Encode())
	return s.httpBase(urlPath, "POST", body, "application/x-www-form-urlencoded", s.options)
}

// PostJSON post json
func (s *EasyHTTP) PostJSON(urlPath string, data interface{}) ([]byte, error) {
	vals, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(vals)
	return s.httpBase(urlPath, "POST", body, "application/json", s.options)
}

func (s *EasyHTTP) httpBase(urlPath, method string, body io.Reader, contentType string, options HTTPOptions) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: options.ConnTimeout,
	}
	req, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "*/*")
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	if len(options.UserAgent) > 0 {
		req.Header.Set("User-Agent", options.UserAgent)
	}
	if len(options.Host) > 0 {
		req.Host = options.Host
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http_status_code: %d, url=%s", resp.StatusCode, urlPath)
	}
	return ioutil.ReadAll(resp.Body)
}

// HTTProxy 转发http请求
//   arg: targetHostSchemed, examples "https://baog.xxxx.com"
func (s *EasyHTTP) HTTProxy(rw http.ResponseWriter, req *http.Request, targetHostSchemed string) {
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
