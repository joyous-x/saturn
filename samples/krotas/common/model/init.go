package model

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"

	"krotas/common/config"
)

const (
	redisDefault    = "default"
	mysqlKeyDefault = "default"
	mysqlKeyMinipro = "minipro"
)

func InitModels() error {
	if nil == dbs.RedisInst(config.GlobalInst().CfgDbs().Redis) {
		return fmt.Errorf("invalid redis instance")
	}
	if nil == dbs.MysqlInst(config.GlobalInst().CfgDbs().Mysql...) {
		return fmt.Errorf("invalid mysql instance")
	}

	if db, err := dbs.MysqlInst().DB(mysqlKeyDefault); err != nil {
		xlog.Error("mysql db:%s error: %v", mysqlKeyDefault, err)
	} else if err := db.Ping(); err != nil {
		xlog.Error("mysql db:%s error: %v", mysqlKeyDefault, err)
	} else {
		xlog.Debug("mysql db:%s ping ok", mysqlKeyDefault)
	}

	redisCon := dbs.RedisInst().Conn(redisDefault)
	defer redisCon.Close()
	if rst, err := redis.String(redisCon.Do("ping")); err != nil {
		xlog.Error("redis ping error: %v", err)
	} else {
		xlog.Debug("redis ping ok: rst=%v", rst)
	}

	xlog.Info("===> Models(%v) init ok", os.Args[0])
	return nil
}
