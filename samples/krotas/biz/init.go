package biz

import (
	"github.com/jinzhu/gorm"
	"github.com/joyous-x/saturn/model/ip2region"
	"github.com/joyous-x/saturn/model/user"
)

func InitSatellates(dbOrm *gorm.DB) {
	ip2region.Init("")
	user.Init(dbOrm)
}
