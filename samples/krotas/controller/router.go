package controller

import (
	"krotas/bizs"
	"krotas/bizs/sv"
	"krotas/bizs/wechat"

	"github.com/gin-gonic/gin"
)

func initRouter(prefix string, r gin.IRouter) {
	wechat.InitRouter("/wx", r)
	r.POST(prefix+"/com/login", bizs.UserLogin)
	r.POST(prefix+"/com/ip2region", bizs.Ip2Region)
	r.POST(prefix+"/com/sms/send", bizs.SendSMS)
	r.POST(prefix+"/com/fdback", bizs.UserFeedback)

	r.POST(prefix+"/sv/parser", sv.URLParser)
	r.POST(prefix+"/sv/vip/types", sv.VipTypes)

	r.POST(prefix+"/tr2cartoon", bizs.Tran2Cartoon)
	r.POST(prefix+"/pay/ali", bizs.AliPay)

}
