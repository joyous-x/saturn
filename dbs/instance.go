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
var gRedises *Redises
var gRedisesOnce sync.Once

// RedisInst 默认使用 useClient = true
func RedisInst(redisConfs ...RedisConf) *Redises {
	if len(redisConfs) > 0 {
		gRedisesOnce.Do(func() {
			redises := &Redises{}
			if err := redises.Init(redisConfs, true); err != nil {
				xlog.Error("init redises err:%v", err)
			} else {
				gRedises = redises
				xlog.Info("===> RedisInst(%v) init ok ", os.Args[0])
			}
		})
	}
	return gRedises
}
