package main

import (
	"github.com/joyous-x/enceladus/common/xlog"
	"github.com/joyous-x/enceladus/ginbox"
	"github.com/gin-gonic/gin"
)

func HttpRouter(r gin.IRouter) error {
	r.GET("/version", version)
	r.POST("/version", version)
	return nil
}

func main() {
	xlog.Debug("-------test start ")

	config := []*ginbox.ServerConfig{
		&ginbox.ServerConfig{
			Name: "",
			Port: 8001,
		},
	}
	_, err := ginbox.InitDefault(config)
	xlog.Debug("-------test InitDefault err: %v ", err)

	HttpRouter(ginbox.DefaultServer())
	ginbox.Default().RunServers()

	xlog.Debug("-------test end ")
}
