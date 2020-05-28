package bizs

import (
	"github.com/joyous-x/saturn/dbs"
	"github.com/joyous-x/saturn/model/ip2region"
	"github.com/joyous-x/saturn/model/user"
)

const (
	mysqlKeyMinipro = "minipro"
)

// Init initialize for all bizs
func Init() {
	ip2region.Init("")

	dbOrm, err := dbs.MysqlInst().DBOrm(mysqlKeyMinipro)
	if err != nil {
		panic("get database fail")
	}
	user.Init(dbOrm)
}
