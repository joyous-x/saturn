package main

import (
	"testing"
	"runtime"
    "path/filepath"
	"github.com/joyous-x/enceladus/common/xlog"
	"github.com/joyous-x/enceladus/config"
)

func TestConfigParser(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	curDir := filepath.Dir(filename)
	xlog.Debug("-------test start ")
	cfgMgr := config.InitGlobalInst("local", curDir)
	xlog.Debug("-------test data: load: %v", cfgMgr)
	cfgLog := cfgMgr.CfgLog()
	xlog.Debug("-------test cfgLog end : %v ", cfgLog)
	cfgDbs := cfgMgr.CfgDbs()
	xlog.Debug("-------test CfgDbs end : %v ", cfgDbs.Redis)
}
