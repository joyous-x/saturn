package main

import (
	"flag"
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/gins"
	"github.com/joyous-x/saturn/dbs"
	
	"krotas/config"
	"krotas/model"
	"krotas/controller"
	"krotas/biz"
)

const (
	env = "local"
	mysqlKeyMinipro = "minipro"
)

func initComponents() {
	// initialize dbs
	if err := model.InitModels(); err != nil {
		panic(err)
	}

	dbOrm, err := dbs.MysqlInst().DBOrm(mysqlKeyMinipro)
	if err != nil {
		panic("init database fail")
	}

	// initialize components
	biz.InitSatellates(dbOrm)
}

func main() {
	xlog.Debug("gins sample ===> start ")

	configPath := flag.String("config", "./env/config/local", "config path")
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
	err := ginbox.Init(cfgMgr.CfgProj().HttpConfs)
	if err != nil {
		xlog.Debug(" ===> ginbox init err: %v ", err)
		return
	} else {
		xlog.Debug(" ===> ginbox init success ")
	}

	// regist routers
	ginbox.HttpRouter(controller.New())
	ginbox.Run()

	xlog.Debug("gins sample ===> end ")
}
