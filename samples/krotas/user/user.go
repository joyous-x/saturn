package user

import (
	"github.com/gin-gonic/gin"
	usercom "github.com/joyous-x/saturn/component/user"
	usermod "github.com/joyous-x/saturn/component/user/model"
)

type userLoginReq struct {
	reqresp.ReqCommon
	usercom.LoginParams
}

type userLoginResp struct {
	reqresp.RespCommon
	usermod.UserInfo
}

// UserLogin user login via phone and third account
// such as wechat, qq, and so on
func UserLogin(c *gin.Context) {

}
