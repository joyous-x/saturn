package user

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(prefix string, r gin.IRouter) {
	r.POST(prefix+"/v1/login", UserLogin)
}
