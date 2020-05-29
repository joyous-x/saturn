package wxcom

import (
	"encoding/json"
	"fmt"

	"github.com/joyous-x/saturn/common/xnet"
)

const (
	urlOauth2WxMiniApp   = "https://api.weixin.qq.com/sns/jscode2session"
	urlOauth2Wx          = "https://api.weixin.qq.com/sns/oauth2/access_token"
	urlAccessTokenPubAcc = "https://api.weixin.qq.com/cgi-bin/token"
)

// AccessTokenInfo access token info
type AccessTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Oauth2Rst wx oauth2 result
type Oauth2Rst struct {
	AccessTokenInfo
	SessionKey   string `json:"session_key"`
	OpenID       string `json:"openid"`
	UnionID      string `json:"unionid"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

// Oauth2WxH5 transfrom code to access_token
func Oauth2WxH5(appID, appSecret, code string) (*Oauth2Rst, error) {
	url := fmt.Sprintf(`%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code`, urlOauth2Wx, appID, appSecret, code)
	return doOauth2(url)
}

// Oauth2WxApp transfrom code to access_token
func Oauth2WxApp(appID, appSecret, code string) (*Oauth2Rst, error) {
	url := fmt.Sprintf(`%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code`, urlOauth2Wx, appID, appSecret, code)
	return doOauth2(url)
}

// Oauth2WxMiniApp authorize client's js_code when miniapp login. if success, we can get openid, etc.
//    url: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func Oauth2WxMiniApp(appID, appSecret, jsCode string) (*Oauth2Rst, error) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", urlOauth2WxMiniApp, appID, appSecret, jsCode)
	return doOauth2(url)
}

// FetchAccessTokenPubAcc get access token for public account
//   url: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func FetchAccessTokenPubAcc(appID, appSecret string) (*AccessTokenInfo, error) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&grant_type=client_credential", urlAccessTokenPubAcc, appID, appSecret)

	client := xnet.EasyHTTP{}
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	accessTokenInfo := &AccessTokenInfo{}
	if err := json.Unmarshal(body, accessTokenInfo); err != nil {
		return accessTokenInfo, err
	}
	return accessTokenInfo, nil
}

func doOauth2(url string) (*Oauth2Rst, error) {
	client := xnet.EasyHTTP{}
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	resp := &Oauth2Rst{}
	if err := json.Unmarshal(body, resp); err != nil {
		return resp, err
	}
	if resp.Errcode != 0 {
		return resp, fmt.Errorf("oauth2 error: code=%v msg=%v", resp.Errcode, resp.Errmsg)
	}
	return resp, nil
}
