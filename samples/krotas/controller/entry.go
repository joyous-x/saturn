package controller

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(prefix string, r gin.IRouter) {
	r.POST(prefix + "/login", UserLogin)

	r.POST(prefix + "/ip2region", Ip2Region)
}
