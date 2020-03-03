package config

import (
	"github.com/joyous-x/saturn/common/xlog"
	"os"
	"path/filepath"
	"sync"
)

func makeConfItem(configPath string) ([]*ConfItem, error) {
	items := []*ConfItem{}

	if len(configPath) < 1 {
		binpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return items, err
		}
		configPath = filepath.Join(binpath, "./env/config/local/")
	}
	dirpath := configPath
	xlog.Debug("makeConfItem dirpath=%v", dirpath)

	items = append(items, &ConfItem{CfgKeyLog, filepath.Join(dirpath, "logs.yaml"), &xlog.XLogConf{}})
	items = append(items, &ConfItem{CfgKeyDbs, filepath.Join(dirpath, "dbs.yaml"), &ConfDbs{}})
	items = append(items, &ConfItem{CfgKeyProj, filepath.Join(dirpath, "proj.yaml"), &ConfProj{}})
	return items, nil
}

var g_configs *ConfigMgr
var g_configs_once sync.Once

// GlobalInst 全局配置管理器实例
func GlobalInst() *ConfigMgr {
	return g_configs
}

// InitGlobalInst 初始化全局配置管理器: args=[configPath ...]
func InitGlobalInst(args ...string) *ConfigMgr {
	g_configs_once.Do(func() {
		configPath := ""
		if len(args) > 0 {
			configPath = args[0]
		}
		xlog.Debug("InitGlobalInst configPath=%v: ready", configPath)

		confs, err := makeConfItem(configPath)
		if err != nil {
			xlog.Error("configmgr makeConfItem err:%v", err)
			return
		}
		configMgr := &ConfigMgr{}
		if err := configMgr.Init(confs); err != nil {
			xlog.Error("configmgr init err:%v", err)
			return
		}
		if err := configMgr.Load(); err != nil {
			xlog.Error("configmgr load err:%v", err)
			return
		}
		g_configs = configMgr
		xlog.Info("===> InitGlobalInst(%v) init ok: %v ", os.Args[0], configPath)
	})
	return g_configs
}
