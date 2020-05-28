package user

import (
	"context"
	"encoding/json"

	"github.com/joyous-x/saturn/foos/user/errors"
	"github.com/joyous-x/saturn/foos/user/model"
)

const (
	// LoginWxMiniApp login by wx app
	LoginWxMiniApp = "wx_miniapp"
	// LoginByWxH5 login by wx h5
	LoginByWxH5 = "wx_h5"
	// LoginByWxQr login by wx qr
	LoginByWxQr = "wx_qr"
	// LoginByWxApp login by wx app
	LoginByWxApp = "wx_app"
	// LoginByMobile login by wx app
	LoginByMobile = "mobile"
)

// LoginParams args
type LoginParams struct {
	InviterUID    string            `json:"inviter_uid"`
	InviteScene   string            `json:"invite_scene"`
	InvitePayload json.RawMessage   `json:"invite_payload,omitempty"`
	Platform      string            `json:"platform"`
	LoginType     string            `json:"login_type"` // 登录方式
	WX            LoginWxParams     `json:"wx"`         // 微信登录
	QQ            LoginQQParams     `json:"qq"`         // QQ登录
	Mobile        LoginMobileParams `json:"mobile"`     // 手机登录
}

// LoginWxParams args
type LoginWxParams struct {
	Code string `json:"code" yaml:"code"`
}

// LoginQQParams args
type LoginQQParams struct {
	AccessToken string `json:"access_token"`
}

// LoginMobileParams args
type LoginMobileParams struct {
	Phone    string `json:"phone"` // 手机号
	Code     string `json:"code"`  // 验证码
	CodeType string `json:"type"`  // 验证码序列号
}

// Login login
func Login(ctx context.Context, req *LoginParams) (*model.UserInfo, error) {
	userInfo := &model.UserInfo{}
	var err error
	switch req.LoginType {
	case LoginByWxQr:
		// 1. appid + appsecret => access_token : (https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183)
		// 2. https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=2 => sdk_ticket
		// 3. sdk_ticket => (client)制作获取 qr 的签名 => (client)获取二维码
	case LoginByWxH5:
		// code(from client) + appid + appsecret => access_token => get informations using api:/sns/xxx
	case LoginByWxApp:
		// code(from client) + appid + appsecret => access_token => get informations using api:/sns/xxx
		userInfo, err = loginByWxApp(ctx, req)
	case LoginByMobile:
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
