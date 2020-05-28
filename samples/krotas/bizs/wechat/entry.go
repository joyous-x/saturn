package wechat

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(prefix string, r gin.IRouter) {
	// authorizate mini applications and get infomations
	r.POST(prefix+"/miniapp/user/login", wxMiniappLogin)
	r.POST(prefix+"/miniapp/user/update", wxMiniappUpdateUser)
	r.POST(prefix+"/miniapp/access_token", wxMiniappAccessToken)

	// wechat public account handler
	r.GET(prefix+"/pubacc/event_message", wxPublicAccountEventMessage)
	r.POST(prefix+"/pubacc/event_message", wxPublicAccountEventMessage)
}
