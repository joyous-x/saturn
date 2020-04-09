package controller

import (
	"github.com/gin-gonic/gin"
	"krotas/biz"
)

func initRouter(prefix string, r gin.IRouter) {
	r.POST(prefix+"/login", biz.UserLogin)
	r.POST(prefix+"/ip2region", biz.Ip2Region)
}
