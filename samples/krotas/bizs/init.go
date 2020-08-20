package bizs

import (
	"github.com/joyous-x/saturn/dbs"
	"github.com/joyous-x/saturn/foos/ip2region"
	"github.com/joyous-x/saturn/foos/user"
)

// Init initialize for all bizs
func Init() {
	if err := InitModels(); err != nil {
		panic(err)
	}

	ip2region.Init("")

	dbOrm, err := dbs.MysqlInst().DBOrm(mysqlKeyMinipro)
	if err != nil {
		panic("get database fail")
	}
	user.Init(dbOrm)
}
