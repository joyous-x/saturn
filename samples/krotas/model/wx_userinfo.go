package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs/jmysql"
	"sync"
	"time"
)

// WxUserInfo 用户基本信息
type WxUserInfo struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	OpenID      string    `json:"openid"`
	UnionID     string    `json:"unionid"`
	SessionKey  string    `json:"session_key"`
	Status      int       `json:"status"`
	NickName    string    `json:"nickname"`
	AvatarURL   string    `json:"avatar_url"`
	Mobile      string    `json:"mobile"`
	InviterID   string    `json:"inviter"`
	Gender      int       `json:"gender"`
	Language    string    `json:"language"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	CreatedTime time.Time `json:"create_time"`
	UpdatedTime time.Time `json:"last_login_time"`
}

// WxUserDao ...
type WxUserDao struct {
	//db *sqlx.DB
	dbOrm *gorm.DB
}

var gWxUserdaoOnce sync.Once
var gWxUserDaoInst *WxUserDao

// WxUserDaoInst ...
func WxUserDaoInst() *WxUserDao {
	gWxUserdaoOnce.Do(func() {
		dbOrm, err := jmysql.GlobalInst().DBOrm(mysqlKeyMinipro)
		if err != nil {
			panic("init database fail")
		}
		gWxUserDaoInst = &WxUserDao{
			dbOrm: dbOrm,
		}
	})
	return gWxUserDaoInst
}

func (w *WxUserDao) tableName(appname string) string {
	return tnUserInfo(appname)
}

// GetUserInfoByUUID 获取用户信息
func (w *WxUserDao) GetUserInfoByUUID(ctx context.Context, appname, uuid string) (*WxUserInfo, error) {
	data := &WxUserInfo{}
	err := error(nil)
	schema := fmt.Sprintf("SELECT * FROM `%s` WHERE `uuid`=?", w.tableName(appname))
	db := w.dbOrm.Raw(schema, uuid).Scan(data)
	if db.Error == sql.ErrNoRows || db.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = db.Error
	}
	return data, err
}

// GetUserInfoByUnionID 获取用户信息
func (w *WxUserDao) GetUserInfoByUnionID(ctx context.Context, appname, unionID string) (*WxUserInfo, error) {
	data := &WxUserInfo{}
	err := error(nil)
	schema := fmt.Sprintf("SELECT * FROM `%s` WHERE `unionid`=?", w.tableName(appname))
	db := w.dbOrm.Raw(schema, unionID).Scan(data)
	if db.Error == sql.ErrNoRows || db.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = db.Error
	}
	return data, err
}

// GetUserInfoByOpenID 获取用户信息
func (w *WxUserDao) GetUserInfoByOpenID(ctx context.Context, appname, openID string) (*WxUserInfo, error) {
	data := &WxUserInfo{}
	err := error(nil)
	schema := fmt.Sprintf("SELECT * FROM `%s` WHERE `openid`=?", w.tableName(appname))
	db := w.dbOrm.Raw(schema, openID).Scan(data)
	if db.Error == sql.ErrNoRows || db.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = db.Error
	}
	return data, err
}

// GetUserInfoByMobile 手机号码获取用户信息
func (w *WxUserDao) GetUserInfoByMobile(ctx context.Context, appname, mobile string) (*WxUserInfo, error) {
	data := &WxUserInfo{}
	err := error(nil)
	schema := fmt.Sprintf("SELECT * FROM `%s` WHERE `mobile`=?", w.tableName(appname))
	db := w.dbOrm.Raw(schema, mobile).Scan(data)
	if db.Error == sql.ErrNoRows || db.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = db.Error
	}
	return data, err
}

// UpdateUserBaseInfo ...
func (w *WxUserDao) UpdateUserBaseInfo(ctx context.Context, appname, uuid, openID, sessionKey string, status int, inviter string) (err error) {
	schema := fmt.Sprintf(`INSERT INTO %s (uuid, openid, session_key, status, inviter, create_time, last_login_time) VALUES (?,?,?,?,?,?,?) 
							ON DUPLICATE KEY UPDATE status=values(status), session_key=values(session_key), last_login_time=values(last_login_time)`, w.tableName(appname))
	db := w.dbOrm.Exec(schema, uuid, openID, sessionKey, status, inviter, time.Now(), time.Now())
	if db.Error != nil {
		xlog.Error("UpdateUserBaseInfo appname=%v uuid=%v inviter=%v error: %v", appname, uuid, inviter, db.Error)
	}
	return db.Error
}

// UpdateUserMobile ...
func (w *WxUserDao) UpdateUserMobile(ctx context.Context, appname, uuid, mobile string) (err error) {
	schema := fmt.Sprintf(`UPDATE %s SET mobile=? WHERE uuid=?`, w.tableName(appname))
	db := w.dbOrm.Exec(schema, mobile, uuid)
	if db.Error != nil {
		xlog.Error("UpdateUserMobile appname=%v uuid=%v mobile=%v error: %v", appname, uuid, mobile, db.Error)
	}
	return db.Error
}

// UpdateUserExtInfo ...
func (w *WxUserDao) UpdateUserExtInfo(ctx context.Context, appname, uuid, unionid, nickname, avatar string, gender int, lang, city, province, country string) (err error) {
	schema := fmt.Sprintf(`UPDATE %s SET unionid=?,nickname=?,avatar_url=?,gender=?,language=?,city=?,province=?,country=? where uuid=?`, w.tableName(appname))
	db := w.dbOrm.Exec(schema, unionid, nickname, avatar, gender, lang, city, province, country, uuid)
	if db.Error != nil {
		xlog.Error("UpdateUserMobile appname=%v uuid=%v nickname=%v error: %v", appname, uuid, nickname, db.Error)
	}
	return db.Error
}
