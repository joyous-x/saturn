package main 

import (
	"benchmysql/mysqltest"
)

func main() {
	db, _ := mysqltest.GetCacheDB()
	mysqltest.Insert2Cache(0, db)
}
