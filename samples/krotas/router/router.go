package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/gins"
	"krotas/user"
	"krotas/wechat"
	"net/http"
)

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
	user.InitRouter("/user", ginbox.Server().Engine())

	return nil
}
