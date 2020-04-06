package controller

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/gins"
	"krotas/wechat"
	"net/http"
)

func New() *Controller {
	return &Controller{} 
}

type Controller struct {
}

func (ctr *Controller) HttpRouter(ginbox *gins.GinBox) error {
	ginbox.Server().Engine().GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Welcome to Saturn"); return })
	ginbox.Server().Engine().POST("/v1/version", ctr.version)
	ginbox.Handle("inner", "POST", "/v1/version", ctr.version)

	ctr.httpRouterStatic(ginbox.Server().Engine())

	wechat.InitRouter("/wx", ginbox.Server().Engine())
	initRouter("/c", ginbox.Server().Engine())

	return nil
}

func (ctr *Controller) version(c *gin.Context) {
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

func (ctr *Controller) httpRouterStatic(r gin.IRouter) error {
	rscPath, err := utils.PathRelative2Bin("./env/rsc/")
	if err != nil {
		return err
	}
	r.Use(static.Serve("/rsc", static.LocalFile(rscPath, false)))
	return nil
}
