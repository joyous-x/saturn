package dbs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joyous-x/saturn/common/xlog"
	// "github.com/jmoiron/sqlx"
)

type MySqls struct {
	dbConns map[string]interface{}
}

func (this *MySqls) Init(dbs []MysqlConf) error {
	this.dbConns = make(map[string]interface{})
	for _, _db := range dbs {
		connMaxLifeTime := 0 * time.Second
		if len(_db.Exts.MaxConnLife) > 0 {
			connMaxLife, err := time.ParseDuration(_db.Exts.MaxConnLife)
			if err != nil {
				xlog.Error("Init ParseDuration failed, key: %v, err: %v", _db.Key, err)
				return err
			}
			connMaxLifeTime = connMaxLife
		}

		if _db.Type == "mysql" {
			m, err := initMySQL(_db.User, _db.Passwd, _db.Host, _db.DbName, _db.Exts.MaxIdle, _db.Exts.MaxOpen, connMaxLifeTime)
			if err != nil {
				xlog.Error("initMySQL failed, key: %v, err: %v", _db.Key, err)
				return err
			}
			this.dbConns[_db.Key] = m
		} else if _db.Type == "mysqlorm" {
			m, err := initMySQLorm(_db.User, _db.Passwd, _db.Host, _db.DbName, _db.Exts.MaxIdle, _db.Exts.MaxOpen, connMaxLifeTime)
			if err != nil {
				xlog.Error("initMySQLorm failed, key: %v, err: %v", _db.Key, err)
				return err
			}
			this.dbConns[_db.Key] = m
		} else {
			xlog.Error("Init MySQL error: invalid type: %v, key: %v", _db.Type, _db.Key)
			return fmt.Errorf("invalid key")
		}
	}
	return nil
}

func (this *MySqls) DB(key string) (*sql.DB, error) {
	if v, ok := this.dbConns[key]; ok {
		return v.(*sql.DB), nil
	} else {
		return nil, fmt.Errorf("Connection not exist")
	}
}

func (this *MySqls) DBOrm(key string) (*gorm.DB, error) {
	if v, ok := this.dbConns[key]; ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, fmt.Errorf("Connection not exist")
	}
}
