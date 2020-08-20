package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/joyous-x/saturn/common/xlog"
)

var g_configs *MgrCenter
var g_configs_once sync.Once

// GlobalInst 全局配置管理器实例
func GlobalInst() *MgrCenter {
	return g_configs
}

// InitGlobalInst 初始化全局配置管理器: args=[configPath ...]
func InitGlobalInst(args ...string) *MgrCenter {
	g_configs_once.Do(func() {
		makeConfigFilePath := func(configDirPath, configFileName string) string {
			if len(configDirPath) < 1 {
				binpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					return ""
				}
				configDirPath = filepath.Join(binpath, "./env/config/local/")
			}
			return filepath.Join(configDirPath, configFileName)
		}

		configPath := func() string {
			if len(args) > 0 {
				return args[0]
			}
			return ""
		}()

		configMgr := &MgrCenter{}
		if err := configMgr.AddConfObjectItem(CfgKeyCommon, makeConfigFilePath(configPath, "config.yaml"), &ComConfig{}); err != nil {
			xlog.Error("configmgr AddConfObjectItem err:%v", err)
			return
		}
		if err := configMgr.Reload(); err != nil {
			xlog.Error("configmgr load err:%v", err)
			return
		}

		g_configs = configMgr
		xlog.Info("===> InitGlobalInst(%v) init ok: %v ", os.Args[0], configPath)
	})
	return g_configs
}
