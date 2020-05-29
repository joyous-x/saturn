package user

import (
	"context"
	"strings"
	"time"

	"github.com/joyous-x/saturn/foos/user/model"
	"github.com/joyous-x/saturn/foos/wechat/wxcom"
)

func getUserInfoByWxApp(openid, accessToken string) (*model.UserInfo, error) {
	userInfo := &model.UserInfo{}
	wxUserInfo, err := wxcom.FetchUserInfo(openid, accessToken)
	if err != nil {
		return userInfo, err
	}

	userInfo.AvatarURL = func(raw string) string {
		// 微信默认返回的是 46*46 大小的头像, 替换为 132*132 的头像
		if raw == "" {
			return raw
		}
		urlTmp := strings.Split(raw, "/")
		return strings.Join(append(urlTmp[0:len(urlTmp)-1], "132"), "/")
	}(wxUserInfo.HeadImgURL)

	userInfo.Gender = func(raw int8) int {
		if raw == 1 || raw == 2 {
			return int(raw)
		}
		return 3
	}(wxUserInfo.Sex)

	userInfo.OpenID = wxUserInfo.OpenID
	userInfo.UnionID = wxUserInfo.UnionID
	userInfo.NickName = wxUserInfo.NickName
	userInfo.CreatedTime = time.Now()
	return userInfo, nil
}

// loginByWxApp register wechat user
func loginByWxApp(ctx context.Context, appid, appname, appsecret string, req *LoginRequest) (*LoginResponse, error) {
	resp := &LoginResponse{}

	oauthRst, err := wxcom.Oauth2WxApp(appid, appsecret, req.WX.Code)
	if err != nil {
		return nil, err
	}

	userInfo, err := getUserInfoByWxApp(oauthRst.OpenID, oauthRst.AccessToken)
	if err != nil {
		return nil, err
	}

	isNewUser, err := updateUserInfo(ctx, appname, userInfo, true)
	if err != nil {
		return nil, err
	}

	resp.NewUser = isNewUser
	resp.Uuid = userInfo.Uuid

	return resp, nil
}

// loginByWxQr login by wechat qr
func loginByWxQr(ctx context.Context, appid, appname, appsecret string, req *LoginRequest) (*LoginResponse, error) {
	resp := &LoginResponse{}

	// TODO

	return resp, nil
}
