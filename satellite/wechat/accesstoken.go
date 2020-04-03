package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"net/http"
	"strings"
)

// AccessTokenResponse access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// FetchPubAcccountUserInfo
func FetchMiniAppAccessToken(appID, appSecret string) (*AccessTokenResponse, error) {
	requestLine := strings.Join([]string{
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=",
		appID, "&secret=", appSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(body, []byte("access_token")) {
		atr := &AccessTokenResponse{}
		err = json.Unmarshal(body, atr)
		if err != nil {
			return nil, err
		}
		return atr, nil
	}

	return nil, fmt.Errorf("%s", "response don't have access_token")
}

var keyAccessToken = func(key string) string {
	return fmt.Sprintf("%s:%s", "wx_ac_token", key)
}

// GetAccessTokenWithCache 获取wechat的access_token
func GetAccessTokenWithCache(conn redis.Conn, appID, appSecret string) (string, error) {
	if nil == conn {
		return "", fmt.Errorf("invalid redis connect")
	}
	defer conn.Close()
	token := getAccessTokenFromRedis(conn, appID, false)
	if len(token) > 0 {
		return token, nil
	}
	resp, err := FetchMiniAppAccessToken(appID, appSecret)
	if err != nil {
		return "", err
	}
	setAccessTokenToRedis(conn, appID, resp.AccessToken, resp.ExpiresIn, false)
	return resp.AccessToken, nil
}

func getAccessTokenFromRedis(conn redis.Conn, appID string, close bool) string {
	if nil == conn {
		return ""
	}
	if close {
		defer conn.Close()
	}
	key := keyAccessToken(appID)
	token, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return ""
	}
	return token
}

func setAccessTokenToRedis(conn redis.Conn, appID string, token string, expiresIn int64, close bool) error {
	if nil == conn {
		return fmt.Errorf("invalid redis conn")
	}
	if close {
		defer conn.Close()
	}
	if expiresIn > 10 {
		expiresIn = expiresIn - 10
	}
	key := keyAccessToken(appID)
	_, err := redis.String(conn.Do("SET", key, token, "ex", expiresIn))
	if err != nil && err != redis.ErrNil {
		return err
	}
	return nil
}
