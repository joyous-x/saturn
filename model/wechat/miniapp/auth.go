package miniapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WxMiniAppAuth 微信登录认证
func WxMiniAppAuth(appID, appSecret, jsCode string) (openID, sessionKey string, err error) {
	URL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, jsCode)

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	wxResp := struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		Errcode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}{}

	err = json.Unmarshal(bodyBytes, &wxResp)
	if err != nil {
		return
	}
	if wxResp.Errcode != 0 {
		err = fmt.Errorf("wxResp:%d,%s", wxResp.Errcode, wxResp.ErrMsg)
		return
	}

	openID = wxResp.OpenID
	sessionKey = wxResp.SessionKey
	return
}
