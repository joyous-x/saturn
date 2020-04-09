package user

import (
	"context"
	"encoding/json"
	"github.com/joyous-x/saturn/model/user/errors"
	"github.com/joyous-x/saturn/model/user/model"
)

const (
	LoginTypeWxApp     = "wx_app"
	LoginTypeWxH5      = "wx_h5"
	LoginTypeWxMiniApp = "wx_miniapp"
	LoginTypeQQ        = "qq"
	LoginTypeMobile    = "mobile"
)

type LoginParams struct {
	InviterUid    string            `json:"inviter_uid"`
	InviteScene   string            `json:"invite_scene"`
	InvitePayload json.RawMessage   `json:"invite_payload,omitempty"`
	Platform      string            `json:"platform"`
	LoginType     string            `json:"login_type"` // 登录方式
	WX            LoginWxParams     `json:"wx"`         // 微信登录
	QQ            LoginQQParams     `json:"qq"`         // QQ登录
	Mobile        LoginMobileParams `json:"mobile"`     // 手机登录
}

type LoginWxParams struct {
	AuthorizationCode string `json:"authorization_code" yaml:"authorization_code"`
}

type LoginQQParams struct {
	AccessToken string `json:"access_token"`
}

type LoginMobileParams struct {
	Phone    string `json:"phone"` // 手机号
	Code     string `json:"code"`  // 验证码
	CodeType string `json:"type"`  // 验证码序列号
}

func Login(ctx context.Context, req *LoginParams) (*model.UserInfo, error) {
	userInfo := &model.UserInfo{}
	var err error
	switch req.LoginType {
	case LoginTypeQQ:
		userInfo, err = loginByQQ(ctx, req)
	case LoginTypeWxH5, LoginTypeWxApp:
		userInfo, err = loginByWX(ctx, req)
	case LoginTypeMobile:
		userInfo, err = loginByMobile(ctx, req)
	default:
		userInfo, err = nil, errors.ErrBadRequest
	}
	if err != nil {
		return userInfo, err
	}
	//> TODO:
	return userInfo, err
}
