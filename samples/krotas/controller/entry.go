package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/gins"
	"github.com/joyous-x/saturn/common/utils"
)

// New create a new Controller
func New() *Controller {
	return &Controller{}
}

// Controller defination for Controller
type Controller struct {
}

// HTTPRouter router for http server
func (ctr *Controller) HTTPRouter(ginbox *gins.GinBox) error {
	ginbox.Server().Engine().POST("/v1/version", ctr.version)
	ginbox.Handle("inner", "POST", "/v1/version", ctr.version)

	ctr.httpRouterStatic(ginbox.Server().Engine())
	initRouter("/api", ginbox.Server().Engine())

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

func (ctr *Controller) httpRouterStatic(r *gin.Engine) error {
	exePath, err := utils.GetExecDirPath()
	if err != nil {
		panic(err.Error())
	}

	r.LoadHTMLGlob(filepath.Join(exePath, "web/templates/*"))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "hello world",
		})
	})
	r.StaticFile("/favicon.ico", filepath.Join(exePath, "web/assert/favicon.ico"))
	r.Use(static.Serve("/web", static.LocalFile(filepath.Join(exePath, "web"), true)))
	r.Use(static.Serve("/rsc", static.LocalFile(filepath.Join(exePath, "rsc"), false)))
	return nil
}
