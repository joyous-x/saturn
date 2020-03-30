package user

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(prefix string, r gin.IRouter) {
	// authorizate mini applications and get infomations
	r.POST(prefix+"/v1/login", UserLogin)
}
