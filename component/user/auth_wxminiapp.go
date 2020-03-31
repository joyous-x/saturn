package user

import (
	"fmt"
	"context"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/wechat"
	"github.com/joyous-x/saturn/component/wechat/miniapp"
	guuid "github.com/google/uuid"
	comerrors "github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/utils"
	"hash/crc64"
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

// LoginByWxMiniApp authorizate and register wechat user
func LoginByWxMiniApp(ctx context.Context, appid, appname, appsecret, jsCode, inviter string) (uuid, token string, isNewUser bool, err error) {
	openID, sessionKey, err := miniapp.WxMiniAppAuth(appid, appsecret, jsCode)
	if err != nil {
		xlog.Error("wxMiniAppLogin appid=%v jscode=%v err=%v", appid, jsCode, err)
		return
	}

	wxUser, err := UserDaoInst().GetUserInfoByOpenID(ctx, appname, openID)
	if err != nil {
		xlog.Error("wxMiniAppLogin GetUserInfoByOpenID appid=%v openID=%v err=%v", appid, openID, err)
		return
	}

	newUUID := func(appname, openid string) string {
		return utils.CalMD5(strings.Replace(guuid.New().String(), "-", "", -1) + appname + openid)
	}
	newToken := func(appname, uuid string) string {
		suffix := strconv.FormatInt(time.Now().UnixNano(), 10) + guuid.New().String()
		hash64 := crc64.Checksum([]byte(appname+uuid+suffix), crc64.MakeTable(crc64.ISO))
		return fmt.Sprintf("%x", hash64)
	}

	if wxUser.Uuid == "" {
		uuid = newUUID(appname, openID)
	} else {
		uuid = wxUser.Uuid
		if wxUser.Status != 0 {
			err = comerrors.ErrAuthForbiden.Err()
			return
		}
	}
	token = newToken(appname, uuid)

	err = UserDaoInst().UpdateUserBaseInfo(ctx, appname, uuid, openID, sessionKey, 0, inviter)
	if err != nil {
		xlog.Error("wxMiniAppLogin UpdateUserBaseInfo appid=%v openid=%v err=%v", appid, openID, err)
		return
	}

	return
}

// updateByWxMiniApp 更新用户信息
func updateByWxMiniApp(ctx context.Context, wxConfig *wechat.WxConfig, reqData *userInfoUpdateReqData) error {
	appname := wxConfig.AppName
	appname, uuid := "", ""

	wxUser, err := UserDaoInst().GetUserInfoByUUID(ctx, appname, uuid)
	if err != nil {
		xlog.Error("GetUserInfoByUUID (%s %s) fail: %v", appname, uuid, err)
		return err
	}

	infos, err := miniapp.DecryptWxUserInfo(reqData.EncryptedData, reqData.Iv, wxUser.SessionKey)
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
