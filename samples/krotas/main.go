package main

import (
	"flag"
	"fmt"

	"github.com/joyous-x/saturn/common/gins"
	"github.com/joyous-x/saturn/common/xlog"

	"krotas/bizs"
	"krotas/bizs/common/config"
	"krotas/controller"
)

const (
	env = "local"
)

func main() {
	xlog.Debug("gins sample ===> start ")

	configPath := flag.String("config", "./config/local", "config path")
	flag.Parse()

	// load configs
	cfgMgr := config.InitGlobalInst(*configPath)
	if cfgMgr == nil {
		panic(fmt.Errorf("invalid config mgr"))
	}

	// make ginbox
	ginbox := gins.DefaultBox()
	if err := ginbox.Init(cfgMgr.GetComConfig().ServerConfs); err != nil {
		panic(fmt.Errorf("ginbox init error: %s", err))
	} else {
		xlog.Debug(" ===> ginbox init success ")
	}

	// initialize models and components
	bizs.Init()
	// http router
	ginbox.HTTPRouter(controller.New())
	// running
	ginbox.Run()

	xlog.Debug("gins sample ===> end ")
}
