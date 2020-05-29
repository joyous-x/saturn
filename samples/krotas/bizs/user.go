package bizs

import (
	"krotas/common/errcode"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/foos/user"
)

// UserLoginReq login request
type UserLoginReq struct {
	reqresp.ReqCommon
	Params *user.LoginRequest `json:"params"`
}

// UserLoginResp login request
type UserLoginResp struct {
	reqresp.RespCommon
	user.LoginResponse
}

// UserLogin user login via phone and third account
// such as wechat, qq, and so on
func UserLogin(c *gin.Context) {
	req := UserLoginReq{}
	ctx, err := reqresp.RequestUnmarshal(c, nil, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, nil)
		return
	}

	info, err := user.Login(ctx, req.Params)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.NewError(errcode.ErrUserLogin.Code, err.Error()), nil)
		return
	}

	resp := &UserLoginResp{
		LoginResponse: *info,
	}
	reqresp.ResponseMarshal(c, errors.OK, resp)
	return
}
