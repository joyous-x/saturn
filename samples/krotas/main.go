package main

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/gins"
	"krotas/router"
)

func main() {
	xlog.Debug("gins sample ===> start ")

	config := []*gins.ServerConfig{
		&gins.ServerConfig{
			Name: "",
			Port: 8001,
		},
	}
	ginbox := gins.DefaultBox()
	err := ginbox.Init(config)
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
