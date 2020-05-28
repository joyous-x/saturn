package bizs

import (
	"krotas/common/errcode"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/common/xnet"
	"github.com/joyous-x/saturn/model/ip2region"
)

var prodMap map[string]string = map[string]string{
	"appid_testA": "secret_testA",
}

type Ip2RegionReq struct {
	reqresp.ReqCommon
	ClientIP string `json:"client_ip"`
	Debug    bool   `json:"debug"`
}

type Ip2RegionResp struct {
	reqresp.RespCommon
	ClientIP string `json:"client_ip"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Province string `json:"province"`
	ISP      string `json:"ISP"`
	CityId   int64  `json:"city_id"`
}

func checkAuth(appid, appSecret, signature string, req *Ip2RegionReq) (bool, error) {
	values := url.Values{}
	values.Add("client_ip", req.ClientIP)
	okSignature := utils.MakeSign(appid, appSecret, values)
	if signature != okSignature {
		xlog.Error("checkAuth proid=%v uniqid=%v curSign=%v okSign=%v", appid, appSecret, signature, okSignature)
		return false, errors.ErrInvalidSign
	}
	return true, nil
}

func Ip2Region(c *gin.Context) {
	req := Ip2RegionReq{}
	resp := Ip2RegionResp{}

	_, err := reqresp.RequestUnmarshal(c, nil, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, &resp)
		return
	}

	if !req.Debug {
		appSecret := ""
		if _, ok := prodMap[req.Common.AppId]; !ok {
			reqresp.ResponseMarshal(c, errors.ErrInvalidAppid, &resp)
			return
		} else {
			appSecret = prodMap[req.Common.AppId]
		}

		signature := c.GetHeader("Authorization")
		if signature == "" {
			reqresp.ResponseMarshal(c, errors.ErrAuthInvalid, &resp)
			return
		}

		if ok, _ := checkAuth(req.Common.AppId, appSecret, signature, &req); !ok {
			reqresp.ResponseMarshal(c, errors.ErrAuthForbiden, &resp)
			return
		}
	}

	clientIP := req.ClientIP
	if len(clientIP) < 1 {
		clientIP = new(xnet.HTTPRealIP).RealIP(c.Request)
	}

	ipInfo, err := ip2region.Inst().MemorySearch(clientIP)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.NewError(errcode.ErrIp2regionMemSearch.Code, err.Error()), &resp)
		return
	}
	resp.ClientIP = clientIP
	resp.Country = ipInfo.Country
	resp.City = ipInfo.City
	resp.Region = ipInfo.Region
	resp.Province = ipInfo.Province
	resp.ISP = ipInfo.ISP
	resp.CityId = ipInfo.CityId
	reqresp.ResponseMarshal(c, errors.OK, &resp)
}
