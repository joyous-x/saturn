package controller

import (
	"github.com/gin-gonic/gin"
	"krotas/bizs"
	"krotas/bizs/sv"
	"krotas/bizs/wechat"
)

func initRouter(prefix string, r gin.IRouter) {
	wechat.InitRouter("/wx", r)
	r.POST(prefix+"/login", bizs.UserLogin)

	r.POST(prefix+"/ip2region", bizs.Ip2Region)
	r.POST(prefix+"/tr2cartoon", bizs.Tran2Cartoon)
	r.POST(prefix+"/pay/ali", bizs.AliPay)
	r.POST(prefix+"/sv/parser", sv.URLParser)
}
