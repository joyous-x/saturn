package user

import (
	"context"

	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/foos/user/model"
	"github.com/joyous-x/saturn/foos/wechat"
	"github.com/joyous-x/saturn/foos/wechat/wxcom"
)

type userInfoUpdateReqData struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

type loginRequestData struct {
	JsCode  string `json:"jscode"`
	Inviter string `json:"inviter"`
}

// LoginByWxMiniApp authorizate and register wechat user
func LoginByWxMiniApp(ctx context.Context, appid, appname, appsecret, jsCode, inviter string) (*LoginResponse, error) {
	resp := &LoginResponse{}

	oauth2Rst, err := wxcom.Oauth2WxMiniApp(appid, appsecret, jsCode)
	if err != nil {
		xlog.Error("LoginByWxMiniApp appid=%v jscode=%v err=%v", appid, jsCode, err)
		return resp, err
	}

	userInfo := &model.UserInfo{
		InviterID:  inviter,
		OpenID:     oauth2Rst.OpenID,
		SessionKey: oauth2Rst.SessionKey,
	}
	isNewUser, err := updateUserInfo(ctx, appname, userInfo, false)
	if err != nil {
		return nil, err
	}

	resp.NewUser = isNewUser
	resp.Uuid = userInfo.Uuid
	return resp, nil
}

// UpdateUserInfoByWxMiniApp 更新用户信息
func UpdateUserInfoByWxMiniApp(ctx context.Context, wxConfig *wechat.WxConfig, reqData *userInfoUpdateReqData) error {
	appname := wxConfig.AppName
	appname, uuid := "", ""

	wxUser, err := UserDaoInst().GetUserInfoByUUID(ctx, appname, uuid)
	if err != nil {
		xlog.Error("GetUserInfoByUUID (%s %s) fail: %v", appname, uuid, err)
		return err
	}

	infos, err := wxcom.DecryptWxUserInfo(reqData.EncryptedData, reqData.Iv, wxUser.SessionKey)
	if err != nil {
		xlog.Error("DecryptWxUserInfo (%s %s) encrypted_data=%v fail: %v", appname, uuid, reqData.EncryptedData, err)
		return err
	}

	err = UserDaoInst().UpdateUserExtInfo(ctx, appname, uuid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
	if err != nil {
		xlog.Error("UpdateUserExtInfo (%s) fail: %v", uuid, err)
	} else {
		xlog.Debug("UpdateUserExtInfo appname=%v uuid=%v nickname=%v avatar=%v", appname, uuid, infos.NickName, infos.AvatarURL)
	}

	return err
}
