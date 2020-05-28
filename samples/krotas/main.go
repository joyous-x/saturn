package main

import (
	"flag"
	"fmt"

	"github.com/joyous-x/saturn/common/gins"
	"github.com/joyous-x/saturn/common/xlog"

	"krotas/bizs"
	"krotas/common/config"
	"krotas/common/model"
	"krotas/controller"
)

const (
	env = "local"
)

func initComponents() {
	// initialize dbs
	if err := model.InitModels(); err != nil {
		panic(err)
	}

	// initialize components
	bizs.Init()
}

func main() {
	xlog.Debug("gins sample ===> start ")

	configPath := flag.String("config", "./config/local", "config path")
	flag.Parse()

	// load configs
	cfgMgr := config.InitGlobalInst(*configPath)
	if cfgMgr == nil {
		panic(fmt.Errorf("invalid config mgr"))
	}

	// initialize models and components
	initComponents()

	// make ginbox
	ginbox := gins.DefaultBox()
	if err := ginbox.Init(cfgMgr.CfgProj().HttpConfs); err == nil {
		xlog.Debug(" ===> ginbox init success ")
		ginbox.HTTPRouter(controller.New())
		ginbox.Run()
	}

	xlog.Debug("gins sample ===> end ")
}
