package dbs

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func initMySQL(user, passwd, host, dbname string, maxIdle, maxActive int, connMaxLife time.Duration) (*sql.DB, error) {
	// data source name : username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8&autocommit=true&loc=%s", user, passwd, host, dbname, "Asia%2FShanghai")
	mysqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mysqldb.SetConnMaxLifetime(connMaxLife)
	mysqldb.SetMaxIdleConns(maxIdle)
	mysqldb.SetMaxOpenConns(maxActive)

	err = mysqldb.Ping()
	if err != nil {
		return mysqldb, err
	}
	return mysqldb, nil
}

func initMySQLorm(user, passwd, host, dbname string, maxIdle, maxActive int, connMaxLife time.Duration) (*gorm.DB, error) {
	// data source name : username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8&autocommit=true&loc=%s", user, passwd, host, dbname, "Asia%2FShanghai")
	mysqlorm, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mysqlorm.DB().SetConnMaxLifetime(connMaxLife)
	mysqlorm.DB().SetMaxIdleConns(maxIdle)
	mysqlorm.DB().SetMaxOpenConns(maxActive)

	err = mysqlorm.DB().Ping()
	if err != nil {
		return mysqlorm, err
	}

	return mysqlorm, nil
}
