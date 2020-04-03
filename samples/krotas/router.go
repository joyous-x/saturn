package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/gins"
	"krotas/controller"
	"krotas/wechat"
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

func httpRouterStatic(r gin.IRouter) error {
	rscPath, err := utils.PathRelative2Bin("./env/rsc/")
	if err != nil {
		return err
	}
	r.Use(static.Serve("/rsc", static.LocalFile(rscPath, false)))
	return nil
}

func HttpRouter(ginbox *gins.GinBox) error {
	ginbox.Server().Engine().GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Welcome to Saturn"); return })
	ginbox.Server().Engine().POST("/v1/version", version)
	ginbox.Handle("inner", "POST", "/v1/version", version)

	httpRouterStatic(ginbox.Server().Engine())

	wechat.InitRouter("/wx", ginbox.Server().Engine())
	controller.InitRouter("/c", ginbox.Server().Engine())

	return nil
}
