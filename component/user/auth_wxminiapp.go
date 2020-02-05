package user

import (
	"context"
	"fmt"	
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/user/errors"
	"github.com/joyous-x/saturn/component/user/model"
	"github.com/joyous-x/saturn/component/wechat"
	"github.com/joyous-x/saturn/component/wechat/miniapp"
	"strconv"
	"strings"
	"time"
)

type userInfoUpdateReqData struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

type loginRequestData struct {
	JsCode  string `json:"jscode"`
	Inviter string `json:"inviter"`
}

type loginResponseData struct {
	UUID      string `json:"uuid"`
	Token     string `json:"token"`
	IsNewUser bool   `json:"is_new_user"`
}

type wxAccessTokenReq struct {
	Appid string `json:"appid"`
}

type wxAccessTokenResp struct {
	Appid string `json:"appid"`
	Token string `json:"access_token"`
}

// loginByWxMiniApp register wechat user
func loginByWxMiniApp(ctx context.Context, wxConfig *wechat.WxConfig, jsCode string) {
	appid := wxConfig.AppID
	appsecret := wxConfig.AppSecret
	appname := wxConfig.AppName

	openID, sessionKey, err := miniapp.WxMiniAppLogin(appid, appsecret, jsCode)
	if err != nil {
		xlog.Error("loginByWxMiniApp appid=%v jscode=%v err=%v", appid, jsCode, err)
		return
	} else {
		xlog.Debug("loginByWxMiniApp appid=%v jscode=%v succ: openid=%v sessionkey=%v", appid, jsCode, openID, sessionKey)
	}

	wxUser, err := model.GetUserInfoByOpenID(ctx, appname, openID)
	if err != nil {
		xlog.Error("loginByWxMiniApp GetUserInfoByOpenID appid=%v openID=%v err=%v", appid, openID, err)
		return
	}

	uuid := wxUser.UUID
	if wxUser.UUID == "" {
		uuid = utils.NewUUID(appname, openID)
		err = model.UpdateUserInfo(ctx, appname, uuid, openID, sessionKey, 0, inviter)
	} else {
		if wxUser.Status != 0 {
			err = errors.ErrAuthForbiden
		}
	}

	if err != nil {
		xlog.Error("wxMiniAppLogin UpdateUserBaseInfo appid=%v openid=%v err=%v", appid, openID, err)
	}

	return err
}

// updateByWxMiniApp 更新用户信息
func updateByWxMiniApp(ctx context.Context, wxConfig *wechat.WxConfig, reqData *userInfoUpdateReqData) error {
	appid := wxConfig.AppID
	appsecret := wxConfig.AppSecret
	appname := wxConfig.AppName

	wxUser, err := model.GetUserInfoByUUID(ctx, appname, uuid)
	if err != nil {
		xlog.Error("GetUserInfoByUUID (%s %s) fail: %v", appname, uuid, err)
		return err
	}

	infos, err := wxminiapp.DecryptWxUserInfo(reqData.EncryptedData, reqData.Iv, wxUser.SessionKey)
	if err != nil {
		xlog.Error("DecryptWxUserInfo (%s %s) encrypted_data=%v fail: %v", appname, uuid, reqData.EncryptedData, err)
		return err
	}

	err = model.UpdateUserExtInfo(ctx, appname, uuid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
	if err != nil {
		xlog.Error("UpdateUserExtInfo (%s) fail: %v", uuid, err)
	} else {
		xlog.Debug("UpdateUserExtInfo appname=%v uuid=%v nickname=%v avatar=%v", appname, uuid, infos.NickName, infos.AvatarURL)
	}

	return err
}

