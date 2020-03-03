package wechat

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(prefix string, r gin.IRouter) {
	r.POST(prefix+"/wx/miniapp/user/login", wxMiniappLogin)
	r.POST(prefix+"/wx/miniapp/user/update", wxMiniappUpdateUser)
	r.POST(prefix+"/wx/miniapp/access_token", wxMiniappAccessToken)
}
