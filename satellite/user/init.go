package user

import (
	"github.com/joyous-x/saturn/satellite/user/model"
	"github.com/jinzhu/gorm"
	"sync"
)

var gUserdaoOnce sync.Once
var gUserDaoInst *model.UserDao

// Init must be called before UserDaoInst
func Init(dbOrm *gorm.DB) error {
	return userDaoInst().SetDbOrm(dbOrm)
}

func userDaoInst() *model.UserDao {
	gUserdaoOnce.Do(func() {
		gUserDaoInst = &model.UserDao{}
	})
	return gUserDaoInst
}

// UserDaoInst ...
func UserDaoInst() *model.UserDao {
	inst := userDaoInst()
	if inst.GetDbOrm() == nil {
		panic("invalid gorm.DB")
	}
	return inst
}