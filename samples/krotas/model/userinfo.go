package model

import (
	"github.com/joyous-x/saturn/component/user/model"
	"github.com/joyous-x/saturn/dbs"
	"sync"
)

var gUserdaoOnce sync.Once
var gUserDaoInst *model.UserDao

// UserDaoInst ...
func UserDaoInst() *model.UserDao {
	gUserdaoOnce.Do(func() {
		dbOrm, err := dbs.MysqlInst().DBOrm(mysqlKeyMinipro)
		if err != nil {
			panic("init database fail")
		}
		gUserDaoInst = &model.UserDao{}
		gUserDaoInst.SetDbOrm(dbOrm)
	})
	return gUserDaoInst
}
