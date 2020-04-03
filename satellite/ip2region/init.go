package ip2region

import (
	"sync"
	"fmt"
	"path/filepath"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
)

//
// Ip2region： 
//     https://github.com/lionsoul2014/ip2region
//

var (
	defaultIp2Region       *Ip2Region
	onceIp2Region         = new(sync.Once)
	defaultDbPath         = ""
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
		dbFilePath := ""
		if len(dbPath) < 1 {
			execPath, err := utils.GetExecPath()
			if err != nil {
				panic(fmt.Sprintf("GetDefaultIp2Region getExecPath error %v", err))
			}
			dbFilePath = filepath.Join(filepath.Dir(execPath), "conf/ip2region.db")
		} else {
			dbFilePath = dbPath
		}
		defaultIp2Region, err := New(dbFilePath)
		if err != nil {
			xlog.Error("GetDefaultIp2Region init error %v", err)
		} else {
			xlog.Debug("GetDefaultIp2Region init config = %v", dbFilePath)
		}
		_, _ = defaultIp2Region.MemorySearch("127.0.0.1")
	})
	return defaultIp2Region
}
