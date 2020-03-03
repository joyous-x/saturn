package biz

import (
	"context"
	"fmt"
	guuid "github.com/google/uuid"
	comerrors "github.com/joyous-x/saturn/common/errors"
	comutils "github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/component/wechat/miniapp"

	"hash/crc64"
	"strconv"
	"strings"
	"time"

	"krotas/config"
	"krotas/model"
)

type userInfoUpdateReqData struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

func MiniAppLogin(ctx context.Context, appname, jsCode, inviter string) (uuid, token string, isNewUser bool, err error) {
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

	wxUser, err := model.UserDaoInst().GetUserInfoByOpenID(ctx, appname, openID)
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

	err = model.PutWxToken(appname, uuid, token)
	if err != nil {
		xlog.Error("wxMiniAppLogin PutWxToken appid=%v openid=%v err=%v", appid, openID, err)
		return
	}

	err = model.UserDaoInst().UpdateUserBaseInfo(ctx, appname, uuid, openID, sessionKey, 0, inviter)
	if err != nil {
		xlog.Error("wxMiniAppLogin UpdateUserBaseInfo appid=%v openid=%v err=%v", appid, openID, err)
		return
	}

	return
}
