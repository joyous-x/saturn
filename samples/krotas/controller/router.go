package controller

import (
	"krotas/bizs"

	"github.com/gin-gonic/gin"
)

func initRouter(prefix string, r gin.IRouter) {
	r.POST(prefix+"/login", bizs.UserLogin)
	r.POST(prefix+"/ip2region", bizs.Ip2Region)
	r.POST(prefix+"/tr2cartoon", bizs.Tran2Cartoon)
	r.POST(prefix+"/pay/ali", bizs.AliPay)
}
