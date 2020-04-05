package biz

import (
	"github.com/jinzhu/gorm"
	"github.com/joyous-x/saturn/satellite/ip2region"
	"github.com/joyous-x/saturn/satellite/user"
)

func InitSatellates(dbOrm *gorm.DB) {
	ip2region.Init("")
	user.Init(dbOrm)
}