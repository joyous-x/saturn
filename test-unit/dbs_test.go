package main

import (
	"testing"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
	"github.com/joyous-x/saturn/dbs/jmysql"
	"github.com/joyous-x/saturn/dbs/jredis"
)

func TestDbs(t *testing.T) {
	// ENV := "local"
	mysqlConf := dbs.MysqlConf{
		Key: "defalut",
		Type: "mysql",
		Host: "10.12.198.188:33061",
		User: "miniuser",
		Passwd:"0sgckpIvpH5s3vmb",
		DbName: "miniprogram",
		Debug: false,
		Exts: dbs.MysqlExts{
			MaxIdle: 10,
		},
	}
	redisConf := dbs.RedisConf{
		Host:"10.12.198.188:63791",
		Passwd:"123.456",
		Dbs: []dbs.RedisDb{
			dbs.RedisDb{Key:"default", Db:0},
		},
		Exts: dbs.RedisExts{
			ConnTimeout: "5s",
		},
	}

	xlog.Debug("-------test start ")
	conn := jredis.GlobalInst(redisConf).Conn("default")
	xlog.Debug("-------test : redisconn=%v ", conn)
	sql, _ := jmysql.GlobalInst(mysqlConf).DB("default")
	xlog.Debug("-------test : sql=%v ", sql)
	sqlOrm, _ := jmysql.GlobalInst().DBOrm("default_orm")
	xlog.Debug("-------test : sqlOrm=%v ", sqlOrm)
	xlog.Debug("-------test end")
}
