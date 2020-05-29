package dao

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/joyous-x/saturn/dbs"
)

const wxTokenExpire = 3600 * 24 * 7

var comconRedisName = "default"
var wxTokenKey = func(appname, uuid string) string {
	return fmt.Sprintf("%s:tk:%s", appname, uuid)
}

// GetWxToken 取回 Token
func GetWxToken(appname, uuid string) (string, error) {
	conn := dbs.RedisInst().Conn(comconRedisName)
	defer conn.Close()
	key := wxTokenKey(appname, uuid)

	token, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return token, err
	}
	return token, err
}

// PutWxToken 存储 token
func PutWxToken(appname, uuid, token string) error {
	conn := dbs.RedisInst().Conn(comconRedisName)
	defer conn.Close()
	key := wxTokenKey(appname, uuid)
	_, err := conn.Do("SETEX", key, wxTokenExpire, token)
	return err
}
