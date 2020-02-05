package wechat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	comerrors "github.com/joyous-x/saturn/common/errors"
	comutils "github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/wechat/miniapp"
	"hash/crc64"
	"strconv"
	"strings"
	"time"

	"krotas/biz"
	"krotas/common"
	"krotas/config"
	"krotas/model"
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

// MiniappWxLogin login a wechat miniprogram
func MiniappWxLogin(c *gin.Context) {
	request := loginRequestData{}
	response := loginResponseData{}

	ctx, appname, _, err := common.RequestUnmarshalNoAuth(c, biz.GetUserInfo, &request)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), response)
		return
	}

	uuid, token, isNewUser, err := wxMiniAppLogin(ctx, appname, request.JsCode, request.Inviter)
	if err != nil {
		common.ResponseMarshal(c, -2, err.Error(), response)
		return
	}

	response.UUID = uuid
	response.Token = token
	response.IsNewUser = isNewUser
	common.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, response)
	return
}

func wxMiniAppLogin(ctx context.Context, appname, jsCode, inviter string) (uuid, token string, isNewUser bool, err error) {
	appInfo, ok := config.GlobalInst().CfgProj().WxApps[appname]
	if !ok {
		err = fmt.Errorf("invalid appname: %v", appname)
		return
	}
	appid := appInfo.AppID
	appsecret := appInfo.AppSecret

	openID, sessionKey, err := miniapp.WxMiniAppAuth(appid, appsecret, jsCode)
	if err != nil {
		xlog.Error("wxMiniAppLogin appid=%v jscode=%v err=%v", appid, jsCode, err)
		return
	} else {
		xlog.Debug("wxMiniAppLogin appid=%v jscode=%v succ: openid=%v sessionkey=%v", appid, jsCode, openID, sessionKey)
	}

	wxUser, err := model.WxUserDaoInst().GetUserInfoByOpenID(ctx, appname, openID)
	if err != nil {
		xlog.Error("wxMiniAppLogin GetUserInfoByOpenID appid=%v openID=%v err=%v", appid, openID, err)
		return
	}

	newUUID := func(appname, openid string) string {
		return comutils.CalMD5(strings.Replace(guuid.New().String(), "-", "", -1) + appname + openid)
	}

	newToken := func(appname, uuid string) string {
		suffix := strconv.FormatInt(time.Now().UnixNano(), 10) + guuid.New().String()
		hash64 := crc64.Checksum([]byte(appname+uuid+suffix), crc64.MakeTable(crc64.ISO))
		return fmt.Sprintf("%x", hash64)
	}

	if wxUser.UUID == "" {
		uuid = newUUID(appname, openID)
	} else {
		uuid = wxUser.UUID
		if wxUser.Status != 0 {
			err = comerrors.ErrAuthForbiden.Err()
			return
		}
	}
	token = newToken(appname, uuid)

	err = model.PutWxToken(appname, uuid, token)
	if err != nil {
		xlog.Error("wxMiniAppLogin PutWxToken appid=%v openid=%v err=%v", appid, openID, err)
		return
	}

	err = model.WxUserDaoInst().UpdateUserBaseInfo(ctx, appname, uuid, openID, sessionKey, 0, inviter)
	if err != nil {
		xlog.Error("wxMiniAppLogin UpdateUserBaseInfo appid=%v openid=%v err=%v", appid, openID, err)
		return
	}

	return
}

// WxUpdateUserInfo 更新用户信息
func WxUpdateUserInfo(c *gin.Context) {
	reqData := userInfoUpdateReqData{}
	ctx, appname, uuid, err := common.RequestUnmarshal(c, biz.GetUserInfo, &reqData)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	_, exist := config.GlobalInst().CfgProj().WxApps[appname]
	if !exist {
		err = fmt.Errorf("invalid appname: %v", appname)
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}

	wxUserInfo, err := model.WxUserDaoInst().GetUserInfoByUUID(ctx, appname, uuid)
	if err != nil {
		xlog.Error("WxUpdateUserInfo GetUserInfoByUUID (%s %s) fail: %v", appname, uuid, err)
		common.ResponseMarshal(c, -2, err.Error(), nil)
	}

	infos, err := miniapp.DecryptWxUserInfo(reqData.EncryptedData, reqData.Iv, wxUserInfo.SessionKey)
	if err != nil {
		xlog.Error("WxUpdateUserInfo DecryptWxUserInfo (%s %s) encrypted_data=%v fail: %v", appname, uuid, reqData.EncryptedData, err)
		common.ResponseMarshal(c, -3, err.Error(), nil)
	}

	err = model.WxUserDaoInst().UpdateUserExtInfo(ctx, appname, uuid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
	if err != nil {
		xlog.Error("WxUpdateUserInfo UpdateUserExtInfo (%s) fail: %v", uuid, err)
		common.ResponseMarshal(c, -4, err.Error(), nil)
	} else {
		xlog.Debug("WxUpdateUserInfo appname=%v uuid=%v nickname=%v avatar=%v", appname, uuid, infos.NickName, infos.AvatarURL)
	}

	common.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, nil)
	return
}