package wechat

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/errors"
	comerrors "github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/wechat"
	"github.com/joyous-x/saturn/component/wechat/miniapp"
	"github.com/joyous-x/saturn/component/wechat/pubacc"
	"github.com/joyous-x/saturn/dbs"
	"github.com/gomodule/redigo/redis"
	"krotas/config"
	"krotas/model"
	"krotas/wechat/biz"
)

type wxMiniappLoginReq struct {
	reqresp.ReqCommon
	JsCode  string `json:"jscode"`
	Inviter string `json:"inviter"`
}

type wxMiniappLoginResp struct {
	reqresp.RespCommon
	UUID      string `json:"uuid"`
	Token     string `json:"token"`
	IsNewUser bool   `json:"is_new_user"`
}

type wxMiniappUpdateUserReq struct {
	reqresp.ReqCommon
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

type wxMiniappAccessTokenReq struct {
	reqresp.ReqCommon
	Appid string `json:"appid"`
}

type wxMiniappAccessTokenResp struct {
	reqresp.RespCommon
	Appid string `json:"appid"`
	Token string `json:"access_token"`
}

// GetUserInfo TODO
func GetUserInfo(appid, uid, authorization string, dataBody []byte) (token string, err error) {
	//> TODO
	return "token_test", nil
}

// wxMiniappLogin login a wechat miniprogram
func wxMiniappLogin(c *gin.Context) {
	req := wxMiniappLoginReq{}
	resp := wxMiniappLoginResp{}

	ctx, err := reqresp.RequestUnmarshal(c, GetUserInfo, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), &resp)
		return
	}

	uuid, token, isNewUser, err := biz.MiniAppLogin(ctx, req.AppId, req.JsCode, req.Inviter)
	if err != nil {
		reqresp.ResponseMarshal(c, -2, err.Error(), &resp)
		return
	}

	resp.UUID = uuid
	resp.Token = token
	resp.IsNewUser = isNewUser
	reqresp.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, &resp)
	return
}

// wUpdateUserInfo 更新用户信息
func wxMiniappUpdateUser(c *gin.Context) {
	req := wxMiniappUpdateUserReq{}
	ctx, err := reqresp.RequestUnmarshal(c, GetUserInfo, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	_, exist := config.GlobalInst().CfgProj().WxApps[req.AppId]
	if !exist {
		err = fmt.Errorf("invalid appid: %v", req.AppId)
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}

	wxUserInfo, err := model.UserDaoInst().GetUserInfoByUUID(ctx, req.AppId, req.Uid)
	if err != nil {
		xlog.Error("WxUpdateUserInfo GetUserInfoByUUID (%s %s) fail: %v", req.AppId, req.Uid, err)
		reqresp.ResponseMarshal(c, -2, err.Error(), nil)
	}

	infos, err := miniapp.DecryptWxUserInfo(req.EncryptedData, req.Iv, wxUserInfo.SessionKey)
	if err != nil {
		xlog.Error("WxUpdateUserInfo DecryptWxUserInfo (%s %s) encrypted_data=%v fail: %v", req.AppId, req.Uid, req.EncryptedData, err)
		reqresp.ResponseMarshal(c, -3, err.Error(), nil)
	}

	err = model.UserDaoInst().UpdateUserExtInfo(ctx, req.AppId, req.Uid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
	if err != nil {
		xlog.Error("WxUpdateUserInfo UpdateUserExtInfo (%s) fail: %v", req.Uid, err)
		reqresp.ResponseMarshal(c, -4, err.Error(), nil)
	} else {
		xlog.Debug("WxUpdateUserInfo appname=%v uuid=%v nickname=%v avatar=%v", req.AppId, req.Uid, infos.NickName, infos.AvatarURL)
	}

	reqresp.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, nil)
	return
}

// wxMiniappAccessToken get a valid access_token for a wechat miniprogram
func wxMiniappAccessToken(c *gin.Context) {
	req := wxMiniappAccessTokenReq{}
	_, err := reqresp.RequestUnmarshal(c, GetUserInfo, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	wxcfg, ok := config.GlobalInst().CfgProj().WxApps[req.AppId]
	if !ok {
		err = fmt.Errorf("invalid appid: %v", req.AppId)
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	token, err := wechat.GetAccessTokenWithCache(dbs.RedisInst().Conn("default"), wxcfg.AppID, wxcfg.AppSecret)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	resp := &wxMiniappAccessTokenResp{
		Appid: wxcfg.AppID,
		Token: token,
	}
	reqresp.ResponseMarshal(c, errors.OK.Code, errors.OK.Msg, resp)
	return
}

// wxPublicAccountEventMessage response for public_account's EventMessage
func wxPublicAccountEventMessage(c *gin.Context) {
	pacfg, ok := config.GlobalInst().CfgProj().WxApps["pubacc"]
	if !ok {
		c.String(http.StatusOK, fmt.Sprintf("invalid pubacc"))
		return
	}
	wxcfg := &wechat.WxConfig{
		AppID:          pacfg.AppID,
		AppSecret:      pacfg.AppSecret,
		EncodingAESKey: pacfg.EncodingAESKey,
		PubAccToken:    pacfg.Token,
	}
	wxcfg.SetRedisFetcher(func() redis.Conn {
		return dbs.RedisInst().Conn("default")
	})

	msgHeader, err := pubacc.ParseRequestFromGin(c, wxcfg, true)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	if msgHeader.EchoStrExist {
		c.String(http.StatusOK, msgHeader.EchoStr)
		return
	}
	replay, err := pubacc.InteractEntry(wxcfg, msgHeader, biz.MyMsgHandler)
	if err == nil && nil != replay {
		c.XML(http.StatusOK, replay)
		return
	}

	retMsg := "success"
	if err != nil {
		retMsg = fmt.Sprintf("%v", err)
		xlog.Error("PublicAccountEventMessage appid=%v error: %v", wxcfg.AppID, err)
	}
	c.String(http.StatusOK, retMsg)
	return
}