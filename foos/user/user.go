package user

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"strconv"
	"strings"
	"time"

	guuid "github.com/google/uuid"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/foos/user/errcode"
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

// LoginRequest args
type LoginRequest struct {
	InviterId     string            `json:"inviter_uuid"`
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

// LoginResponse ...
type LoginResponse struct {
	Uuid    string `json:"uuid"`
	NewUser int    `json:"new_user"`
	Token   string `json:"token"`
}

// Login login
func Login(ctx context.Context, req *LoginRequest) (resp *LoginResponse, err error) {
	appid, appname, appsecret := "", "", ""

	switch req.LoginType {
	case LoginByWxQr:
		// 1. appid + appsecret => access_token : (https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183)
		// 2. https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=2 => sdk_ticket
		// 3. sdk_ticket => (client)制作获取 qr 的签名 => (client)获取二维码
		resp, err = loginByWxQr(ctx, appid, appname, appsecret, req)
	case LoginByWxH5:
		// code(from client) + appid + appsecret => access_token => get informations using api:/sns/xxx
		fallthrough
	case LoginByWxApp:
		// code(from client) + appid + appsecret => access_token => get informations using api:/sns/xxx
		resp, err = loginByWxApp(ctx, appid, appname, appsecret, req)
	case LoginByMobile:
		resp, err = loginByMobile(ctx, appname, req)
	default:
		resp, err = nil, errors.ErrBadRequest
	}
	return
}

func newToken(appname, uuid string) string {
	suffix := strconv.FormatInt(time.Now().UnixNano(), 10) + guuid.New().String()
	hash64 := crc64.Checksum([]byte(appname+uuid+suffix), crc64.MakeTable(crc64.ISO))
	return fmt.Sprintf("%x", hash64)
}

func updateUserInfo(ctx context.Context, appname string, infos *model.UserInfo, updateExtInfo bool) (int, error) {
	wxUser, err := UserDaoInst().GetUserInfoByOpenID(ctx, appname, infos.OpenID)
	if err != nil {
		xlog.Error("UserDaoInst().GetUserInfoByOpenID appname=%v openID=%v err=%v", appname, infos.OpenID, err)
		return 0, err
	}

	newUUID := func(appname, openid string) string {
		return utils.CalMD5(strings.Replace(guuid.New().String(), "-", "", -1) + appname + openid)
	}

	isNewUser := 0
	if wxUser.Uuid == "" {
		infos.Uuid = newUUID(appname, infos.OpenID)
		isNewUser = 1
	} else {
		if wxUser.Status != 0 {
			return isNewUser, errcode.ErrAuthForbiden
		}
		if infos.Uuid != wxUser.Uuid {
			return isNewUser, errcode.ErrConfusedUuid
		}
	}

	err = UserDaoInst().UpdateUserBaseInfo(ctx, appname, infos.Uuid, infos.OpenID, infos.SessionKey, 0, infos.InviterID)
	if err != nil {
		return isNewUser, err
	}

	if updateExtInfo {
		err = UserDaoInst().UpdateUserExtInfo(ctx, appname, infos.Uuid, infos.UnionID, infos.NickName, infos.AvatarURL, infos.Gender, infos.Language, infos.City, infos.Province, infos.Country)
		if err != nil {
			return isNewUser, err
		}
	}

	return isNewUser, nil
}
