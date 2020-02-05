package wechat


import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/dbs/jredis"
	"github.com/joyous-x/saturn/component/wechat"
	"krotas/common"
	"krotas/config"
)


type wxAccessTokenReq struct {
	Appid string `json:"appid"`
}

type wxAccessTokenResp struct {
	Appid string `json:"appid"`
	Token string `json:"access_token"`
}


// MiniappWxAccessToken get a valid access_token for a wechat miniprogram
func MiniappWxAccessToken(c *gin.Context) {
	reqData := wxAccessTokenReq{}
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
	token, err := wechat.GetAccessTokenWithCache(jredis.GlobalInst().Conn("default"), wxcfg.AppID, wxcfg.AppSecret)
	if err != nil {
		common.ResponseMarshal(c, -1, err.Error(), nil)
		return
	}
	respData := &wxAccessTokenResp{
		Appid: wxcfg.AppID,
		Token: token,
	}
	common.ResponseMarshal(c, comerrors.OK.Code, comerrors.OK.Msg, respData)
	return
}
