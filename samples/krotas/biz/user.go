package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/errors"
	usercom "github.com/joyous-x/saturn/model/user"
	usermod "github.com/joyous-x/saturn/model/user/model"
	kerrs "krotas/errors"
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
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, nil)
		return
	}

	info, err := usercom.Login(ctx, req.Params)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.NewError(kerrs.ErrUserLogin.Code, err.Error()), nil)
		return
	}

	resp := &userLoginResp{
		User: info,
	}
	reqresp.ResponseMarshal(c, errors.OK, resp)
	return
}
