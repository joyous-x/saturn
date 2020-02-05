package pubacc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PubAccountUserInfo public account用户信息定义
//   UserInfoResponse https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140839 字段说明
type PubAccountUserInfo struct {
	Subscribe      int8   `json:"subscribe"` //是否订阅该公众号标识
	OpenID         string `json:"openid"`    //用户的标识，对当前公众号唯一
	UnionID        string `json:"unionid"`
	NickName       string `json:"nickname"`
	City           string `json:"city"`
	Country        string `json:"country"`
	Province       string `json:"province"`
	Sex            int8   `json:"sex"`
	Language       string `json:"language"`
	HeadImgUrl     string `json:"head_img_url"`
	SubscribeTime  int64  `json:"subscribe_time"`
	Remark         string `json:"remark"`
	GroupID        int    `json:"groupid"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

// FetchPubAcccountUserInfo 获取public account用户信息
//   获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err
func FetchPubAcccountUserInfo(openid string, AccessToken string) (*PubAccountUserInfo, error) {
	requestLine := strings.Join([]string{
		"https://api.weixin.qq.com/cgi-bin/user/info?lang=zh_CN&access_token=",
		AccessToken, "&openid=", openid}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	userInfo := &PubAccountUserInfo{}
	err = json.Unmarshal(body, userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.OpenID == "" {
		return nil, fmt.Errorf("返回数据json为异常结构, req:%s; body:%s", requestLine, string(body))
	}

	return userInfo, err
}
