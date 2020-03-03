package wechat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	comerrors "github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/wechat"
	"github.com/joyous-x/saturn/component/wechat/miniapp"
	"github.com/joyous-x/saturn/dbs"
	"krotas/common"
	"krotas/config"
	"krotas/model"
	"krotas/wechat/biz"
)

type wxMiniappLoginReq struct {
	JsCode  string `json:"jscode"`
	Inviter string `json:"inviter"`
}

type wxMiniappLoginResp struct {
	UUID      string `json:"uuid"`
	Token     string `json:"token"`
	IsNewUser bool   `json:"is_new_user"`
}

type wxMiniappUpdateUserReq struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

type wxMiniappAccessTokenReq struct {
	Appid string `json:"appid"`
}

type wxMiniappAccessTokenResp struct {
	Appid string `json:"appid"`
	Token string `json:"access_token"`
}

// GetUserInfo TODO
func GetUserInfo(appname, uuid string) (token string, err error) {
	return "token_test", nil
}

// wxMiniappLogin login a wechat miniprogram
func wxMiniappLogin(c *gin.Context) {
	request := wxMiniappLoginReq{}
	response := wxMiniappLoginResp{}

	ctx, appname, _, err := common.RequestUnmarshalNoAuth(c, GetUserInfo, &request)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), response)
		return
	}

	uuid, token, isNewUser, err := biz.MiniAppLogin(ctx, appname, request.JsCode, request.Inviter)
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

// wUpdateUserInfo 更新用户信息
func wxMiniappUpdateUser(c *gin.Context) {
	reqData := wxMiniappUpdateUserReq{}
	ctx, appname, uuid, err := common.RequestUnmarshal(c, GetUserInfo, &reqData)
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

	wxUserInfo, err := model.UserDaoInst().GetUserInfoByUUID(ctx, appname, uuid)
	if err != nil {
		xlog.Error("WxUpdateUserInfo GetUserInfoByUUID (%s %s) fail: %v", appname, uuid, err)
		common.ResponseMarshal(c, -2, err.Error(), nil)
	}

	infos, err := miniapp.DecryptWxUserInfo(reqData.EncryptedData, reqData.Iv, wxUserInfo.SessionKey)
	if err != nil {
		xlog.Error("WxUpdateUserInfo DecryptWxUserInfo (%s %s) encrypted_data=%v fail: %v", appname, uuid, reqData.EncryptedData, err)
		common.ResponseMarshal(c, -3, err.Error(), nil)
	}

	err = model.UserDaoInst().UpdateUserExtInfo(ctx, appname, uuid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
	if err != nil {
		xlog.Error("WxUpdateUserInfo UpdateUserExtInfo (%s) fail: %v", uuid, err)
		common.ResponseMarshal(c, -4, err.Error(), nil)
	} else {
		xlog.Debug("WxUpdateUserInfo appname=%v uuid=%v nickname=%v avatar=%v", appname, uuid, infos.NickName, infos.AvatarURL)
	}

	common.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, nil)
	return
}

// wxMiniappAccessToken get a valid access_token for a wechat miniprogram
func wxMiniappAccessToken(c *gin.Context) {
	reqData := wxMiniappAccessTokenReq{}
	_, appname, _, err := common.RequestUnmarshal(c, GetUserInfo, &reqData)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	wxcfg, ok := config.GlobalInst().CfgProj().WxApps[appname]
	if !ok {
		err = fmt.Errorf("invalid appname: %v", appname)
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	token, err := wechat.GetAccessTokenWithCache(dbs.RedisInst().Conn("default"), wxcfg.AppID, wxcfg.AppSecret)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	respData := &wxMiniappAccessTokenResp{
		Appid: wxcfg.AppID,
		Token: token,
	}
	common.ResponseMarshal(c, errors.OK.Code, errors.OK.Msg, respData)
	return
}
