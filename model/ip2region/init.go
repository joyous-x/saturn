package ip2region

import (
	"fmt"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"path/filepath"
	"sync"
)

//
// Ip2region：
//     https://github.com/lionsoul2014/ip2region
//

var (
	defaultIp2Region *Ip2Region
	onceIp2Region    = new(sync.Once)
	defaultDbPath    = ""
)

// Init 如果有初始化需求可以在这里定义
func Init(dbPath string) {
	defaultIp2RegionInst(dbPath)
}

func Inst() *Ip2Region {
	return defaultIp2Region
}

func defaultIp2RegionInst(dbPath string) *Ip2Region {
	onceIp2Region.Do(func() {
		dbFilePath := dbPath
		if len(dbFilePath) < 1 {
			execDirPath, err := utils.GetExecDirPath()
			if err != nil {
				panic(fmt.Sprintf("defaultIp2RegionInst getExecPath error %v", err))
			}
			dbFilePath = filepath.Join(execDirPath, "config/ip2region/ip2region.db")
		}
		tmpIp2Region, err := New(dbFilePath)
		if err == nil {
			_, err = tmpIp2Region.MemorySearch("127.0.0.1")
		}
		if err != nil {
			panic(fmt.Sprintf("new Ip2Region error: %v", err))
		}
		defaultIp2Region = tmpIp2Region
		xlog.Debug("defaultIp2RegionInst init ok: %v", dbFilePath)
	})
	return defaultIp2Region
}
