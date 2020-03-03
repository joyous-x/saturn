package main

import (
	"fmt"
	"flag"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/gins"
	"krotas/config"
	"krotas/router"
	"krotas/model"
)

const (
	env = "local"
)

func main() {
	xlog.Debug("gins sample ===> start ")

	configPath := flag.String("config", "./env/config/local", "config path")
	flag.Parse()

	cfgMgr := config.InitGlobalInst(*configPath)
	if cfgMgr == nil {
		panic(fmt.Errorf("invalid config mgr"))
	}

	if err := model.InitModels(); err != nil {
		panic(err)
	}

	ginbox := gins.DefaultBox()
	err := ginbox.Init(cfgMgr.CfgProj().HttpConfs)
	if err != nil {
		xlog.Debug(" ===> ginbox init err: %v ", err)
		return
	} else {
		xlog.Debug(" ===> ginbox init success ")
	}

	router.HttpRouter(ginbox)
	ginbox.Run()

	xlog.Debug("gins sample ===> end ")
}
