package jmysql

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
	"sync"
)

var g_mySqls *MySqls
var g_mysqls_once sync.Once

func GlobalInst(mysqlConfs ...dbs.MysqlConf) *MySqls {
	if len(mysqlConfs) > 0 {
		g_mysqls_once.Do(func() {
			mysqlConfItems := mysqlConfs
			obj := &MySqls{}
			if err := obj.Init(mysqlConfItems); err != nil {
				xlog.Error("init MySqls err:%v", err)
			} else {
				g_mySqls = obj
			}
		})
	}
	return g_mySqls
}
