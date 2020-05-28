package common

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/joyous-x/saturn/foos/wechat"
)

var keyAccessToken = func(key string) string {
	return fmt.Sprintf("%s:%s", "wx_ac_token", key)
}

// GetPubAccAccessToken 获取access_token，并维持缓存
func GetPubAccAccessToken(s *wechat.WxConfig) string {
	token, err := getPubAccAccessToken(s.RedisFetcher(), s.AppID, s.AppSecret)
	if err != nil {
		return ""
	}
	return token
}

func getPubAccAccessToken(conn redis.Conn, appID, appSecret string) (string, error) {
	if nil == conn {
		return "", fmt.Errorf("invalid redis connect")
	}
	defer conn.Close()
	token := getAccessTokenFromRedis(conn, appID, false)
	if len(token) > 0 {
		return token, nil
	}
	resp, err := wechat.FetchAccessTokenPubAcc(appID, appSecret)
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
