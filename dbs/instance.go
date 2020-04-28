package dbs

import (
	"os"
	"sync"

	"github.com/joyous-x/saturn/common/xlog"
)

//-----------------------------------------

var gMySqls *MySqls
var gMysqlsOnce sync.Once

// MysqlInst ...
func MysqlInst(mysqlConfs ...MysqlConf) *MySqls {
	if len(mysqlConfs) > 0 {
		gMysqlsOnce.Do(func() {
			mysqlConfItems := mysqlConfs
			obj := &MySqls{}
			if err := obj.Init(mysqlConfItems); err != nil {
				xlog.Error("init MySqls err:%v", err)
			} else {
				gMySqls = obj
				xlog.Info("===> MysqlInst(%v) init ok ", os.Args[0])
			}
		})
	}
	return gMySqls
}

//-----------------------------------------
var gRedisPools *RedisPools
var gRedisPoolsOnce sync.Once

// RedisInst ...
func RedisInst(redisConfs ...RedisConf) *RedisPools {
	if len(redisConfs) > 0 {
		gRedisPoolsOnce.Do(func() {
			redisConfData := redisConfs[0]
			redisPools := &RedisPools{}
			if err := redisPools.Init(&redisConfData); err != nil {
				xlog.Error("init redispools err:%v", err)
			} else {
				gRedisPools = redisPools
				xlog.Info("===> RedisInst(%v) init ok ", os.Args[0])
			}
		})
	}
	return gRedisPools
}

//-----------------------------------------
var gRedisClients *RedisClients
var gRedisClientsOnce sync.Once

// RedisInstEx ...
func RedisInstEx(redisConfs ...RedisConf) *RedisClients {
	if len(redisConfs) > 0 {
		gRedisClientsOnce.Do(func() {
			redisConfData := redisConfs[0]
			redisClients := &RedisClients{}
			if err := redisClients.Init(&redisConfData); err != nil {
				xlog.Error("init redisClients err:%v", err)
			} else {
				gRedisClients = redisClients
				xlog.Info("===> RedisInstEx(%v) init ok ", os.Args[0])
			}
		})
	}
	return gRedisClients
}
