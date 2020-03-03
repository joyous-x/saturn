package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/gins"
	"github.com/joyous-x/saturn/common/utils"
	"krotas/biz"
	wxbiz "krotas/biz/wechat"
)

func HttpRouter(ginbox *gins.GinBox) error {
	ginbox.Server().Engine().POST("/v1/version", version)
	ginbox.Handle("inner", "POST", "/v1/ex/version", version)

	httpRouterStatic(ginbox.Server().Engine())

	ginbox.Server().Engine().POST("/v1/miniapp/wx/login", wxbiz.MiniappWxLogin)
	ginbox.Server().Engine().POST("/v1/miniapp/wx/access_token", wxbiz.MiniappWxAccessToken)

	ginbox.Server().Engine().POST("/v1/user/login", biz.UserLogin)
	return nil
}

func httpRouterStatic(r gin.IRouter) error {
	rscPath, err := utils.PathRelative2Bin("./env/rsc/")
	if err != nil {
		return err
	}
	r.Use(static.Serve("/rsc", static.LocalFile(rscPath, false)))
	return nil
}
