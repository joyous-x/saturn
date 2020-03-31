package user

import (
	"github.com/joyous-x/saturn/component/user/model"
	"github.com/jinzhu/gorm"
	"sync"
)

var gUserdaoOnce sync.Once
var gUserDaoInst *model.UserDao

// Init must be called before UserDaoInst
func Init(dbOrm *gorm.DB) error {
	return UserDaoInst().SetDbOrm(dbOrm)
}

// UserDaoInst ...
func UserDaoInst() *model.UserDao {
	gUserdaoOnce.Do(func() {
		gUserDaoInst = &model.UserDao{}
	})
	if gUserDaoInst.GetDbOrm() == nil {
		panic("invalid gorm.DB")
	}
	return gUserDaoInst
}