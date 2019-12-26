package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/gins"
	"net/http"
)

func version(c *gin.Context) {
	datas := map[string]string{
		"data":   "hello world",
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}
	if "GET" == c.Request.Method {
		datas["args"] = c.Query("args")
	} else if "POST" == c.Request.Method {
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, datas)
}

func HttpRouter(ginbox *gins.GinBox) error {
	ginbox.Handle("default", "POST", "/v1", version)
	ginbox.Server().Engine().POST("/v2", version)
	return nil
}

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

	HttpRouter(ginbox)
	ginbox.Run()

	xlog.Debug("gins sample ===> end ")
}
