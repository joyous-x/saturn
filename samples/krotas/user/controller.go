package user

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/errors"
	usercom "github.com/joyous-x/saturn/component/user"
	usermod "github.com/joyous-x/saturn/component/user/model"
)

type userLoginReq struct {
	reqresp.ReqCommon
	Params  *usercom.LoginParams `json:"params"`
}

type userLoginResp struct {
	reqresp.RespCommon
	User     *usermod.UserInfo `json:"user"`
}

// UserLogin user login via phone and third account
// such as wechat, qq, and so on
func UserLogin(c *gin.Context) {
	req := userLoginReq{}
	ctx, err := reqresp.RequestUnmarshal(c, nil, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}

	info, err := usercom.Login(ctx, req.Params)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}

	resp := &userLoginResp{
		User: info,
	}
	reqresp.ResponseMarshal(c, errors.OK.Code, errors.OK.Msg, resp)
	return
}
