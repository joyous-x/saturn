package wxcom

import (
	"encoding/json"
	"fmt"

	"github.com/joyous-x/saturn/common/xnet"
)

const (
	apiUserInfoDefault = "https://api.weixin.qq.com/sns/userinfo"
	apiUserInfoPubAcc  = "https://api.weixin.qq.com/cgi-bin/user/info"
)

// PubAccountUserInfo public account用户信息定义
//   UserInfoResponse https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140839 字段说明
type PubAccountUserInfo struct {
	CommonUserInfo
	Subscribe      int8   `json:"subscribe"` //是否订阅该公众号标识
	Language       string `json:"language"`
	SubscribeTime  int64  `json:"subscribe_time"`
	Remark         string `json:"remark"`
	GroupID        int    `json:"groupid"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

// CommonUserInfo user info
type CommonUserInfo struct {
	OpenID     string   `json:"openid"` //用户的标识，对当前公众号唯一
	NickName   string   `json:"nickname"`
	Sex        int8     `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"head_img_url"`
	UnionID    string   `json:"unionid"`
	Privilege  []string `json:"privilege"`
}

// FetchUserInfo 获取用户信息
//   url: https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
func FetchUserInfo(openid, accessToken string) (*CommonUserInfo, error) {
	userInfo := &CommonUserInfo{}
	if err := fetchUserInfo(apiUserInfoDefault, accessToken, openid, userInfo); err != nil {
		return userInfo, err
	}
	if userInfo.OpenID == "" {
		return userInfo, fmt.Errorf("invalid user info")
	}
	return userInfo, nil
}

// FetchPubAccUserInfo 获取public account用户信息
//   获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err
//   url: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html
func FetchPubAccUserInfo(openid, accessToken string) (*PubAccountUserInfo, error) {
	userInfo := &PubAccountUserInfo{}
	if err := fetchUserInfo(apiUserInfoPubAcc, accessToken, openid, userInfo); err != nil {
		return userInfo, err
	}
	if userInfo.OpenID == "" {
		return userInfo, fmt.Errorf("invalid user info")
	}
	return userInfo, nil
}

func fetchUserInfo(host, openid, accessToken string, userInfo interface{}) error {
	url := fmt.Sprintf(`%s?access_token=%s&openid=%s&lang=zh_CN`, host, accessToken, openid)

	client := xnet.EasyHTTP{}
	body, err := client.Get(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, userInfo)
}
