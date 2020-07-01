package bizs

import (
	"fmt"
	"os"

	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"

	"krotas/bizs/common/config"
)

const (
	redisDefault    = "default"
	mysqlKeyDefault = "default"
	mysqlKeyMinipro = "minipro"
)

func InitModels() error {
	if nil == dbs.RedisInst(config.GlobalInst().GetComConfig().Redis...) {
		return fmt.Errorf("invalid redis instance")
	}
	if nil == dbs.MysqlInst(config.GlobalInst().GetComConfig().Mysql...) {
		return fmt.Errorf("invalid mysql instance")
	}

	if err := dbs.MysqlInst().Ping(mysqlKeyDefault); err != nil {
		xlog.Error("mysql db:%s error: %v", mysqlKeyDefault, err)
	} else {
		xlog.Debug("mysql db:%s ping ok", mysqlKeyDefault)
	}

	if rst, err := dbs.RedisInst().Ping(redisDefault); err != nil {
		xlog.Error("redis ping error: %v", err)
	} else {
		xlog.Debug("redis ping ok: rst=%v", rst)
	}

	xlog.Info("===> Models(%v) init ok", os.Args[0])
	return nil
}
