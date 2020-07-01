package dbs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joyous-x/saturn/common/xlog"

	_ "github.com/go-sql-driver/mysql"
)

// MySqls ...
type MySqls struct {
	dbConns map[string]interface{}
}

// Init ...
func (m *MySqls) Init(dbs []MysqlConf) error {
	m.dbConns = make(map[string]interface{})
	for _, _db := range dbs {
		connMaxLifeTime := 0 * time.Second
		if len(_db.Exts.MaxConnLife) > 0 {
			connMaxLife, err := time.ParseDuration(_db.Exts.MaxConnLife)
			if err != nil {
				xlog.Error("Init ParseDuration failed, name: %v, err: %v", _db.Name, err)
				return err
			}
			connMaxLifeTime = connMaxLife
		}

		if _db.Type == "mysql" {
			dbTmp, err := initMySQL(_db.Dsn, _db.Exts.MaxIdle, _db.Exts.MaxOpen, connMaxLifeTime)
			if err != nil {
				xlog.Error("initMySQL failed, name: %v, err: %v", _db.Name, err)
				return err
			}
			m.dbConns[_db.Name] = dbTmp
		} else if _db.Type == "mysqlorm" {
			dbTmp, err := initMySQLorm(_db.Dsn, _db.Exts.MaxIdle, _db.Exts.MaxOpen, connMaxLifeTime)
			if err != nil {
				xlog.Error("initMySQLorm failed, name: %v, err: %v", _db.Name, err)
				return err
			}
			m.dbConns[_db.Name] = dbTmp
		} else {
			xlog.Error("Init MySQL error: invalid type: %v, name: %v", _db.Type, _db.Name)
			return fmt.Errorf("invalid name")
		}
	}
	return nil
}

// Ping ...
func (m *MySqls) Ping(key string) error {
	if db, err := m.DB(key); err == nil {
		return db.Ping()
	}
	if db, err := m.DBOrm(key); err == nil && db.DB() != nil {
		return db.DB().Ping()
	}
	return fmt.Errorf("invalid key(%s)", key)
}

// DB ...
func (m *MySqls) DB(key string) (*sql.DB, error) {
	if v, ok := m.dbConns[key]; ok {
		return v.(*sql.DB), nil
	} else {
		return nil, fmt.Errorf("Connection not exist")
	}
}

// DBOrm ...
func (m *MySqls) DBOrm(key string) (*gorm.DB, error) {
	if v, ok := m.dbConns[key]; ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, fmt.Errorf("Connection not exist")
	}
}

func initMySQL(dsn string, maxIdle, maxActive int, connMaxLife time.Duration) (*sql.DB, error) {
	// data source name : username:password@protocol(address)/dbname?param=value
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

func initMySQLorm(dsn string, maxIdle, maxActive int, connMaxLife time.Duration) (*gorm.DB, error) {
	// data source name : username:password@protocol(address)/dbname?param=value
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
