package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
