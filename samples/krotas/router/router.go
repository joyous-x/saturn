package router

import(
	"github.com/joyous-x/saturn/gins"
	"krotas/biz"
)

func HttpRouter(ginbox *gins.GinBox) error {
	ginbox.Server().Engine().POST("/v1/version", version)
	ginbox.Handle("default", "POST", "/v1/ex/version", version)

	ginbox.Server().Engine().POST("/v1/miniapp/wx/login", biz.MiniappWxLogin)
	ginbox.Server().Engine().POST("/v1/miniapp/wx/access_token", biz.MiniappWxAccessToken)

	ginbox.Server().Engine().POST("/v1/user/login", biz.UserLogin)
	return nil
}
