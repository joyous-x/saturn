package dbs

import (
	"github.com/joyous-x/saturn/common/xlog"
	"sync"
	"os"
)

//-----------------------------------------

var g_mySqls *MySqls
var g_mysqls_once sync.Once

func MysqlInst(mysqlConfs ...MysqlConf) *MySqls {
	if len(mysqlConfs) > 0 {
		g_mysqls_once.Do(func() {
			mysqlConfItems := mysqlConfs
			obj := &MySqls{}
			if err := obj.Init(mysqlConfItems); err != nil {
				xlog.Error("init MySqls err:%v", err)
			} else {
				g_mySqls = obj
				xlog.Info("===> MysqlInst(%v) init ok ", os.Args[0])
			}
		})
	}
	return g_mySqls
}

//-----------------------------------------
var g_redisPools *RedisPools
var g_redis_once sync.Once

func RedisInst(redisConfs ...RedisConf) *RedisPools {
	if len(redisConfs) > 0 {
		g_redis_once.Do(func() {
			redisConfData := redisConfs[0]
			redisPools := &RedisPools{}
			if err := redisPools.Init(&redisConfData); err != nil {
				xlog.Error("init redispools err:%v", err)
			} else {
				g_redisPools = redisPools
				xlog.Info("===> RedisInst(%v) init ok ", os.Args[0])
			}
		})
	}
	return g_redisPools
}
