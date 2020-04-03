package controller

import (
	"net/url"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/satellite/ip2region"
)

type Ip2RegionReq struct {
	reqresp.ReqCommon
	ClientIP    string `json:"client_ip"`
}

type Ip2RegionResp struct {
	reqresp.RespCommon
	ClientIP   string `json:"client_ip"`
	Country string `json:"country"`
	City string `json:"city"`
	Region string `json:"region"`
	Province string `json:"province"`
	ISP string `json:"ISP"`
	CityId  int64 `json:"city_id"`
}

func checkAuth(appid, appSecret, signature string, req *Ip2RegionReq) (bool,error) {
	values := url.Values{
		//> TODO: req ---> values
	}	
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
	appSecret := "" //> TODO:

	_, err := reqresp.RequestUnmarshal(c, nil, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, -1, err.Error(), &resp)
		return
	}

	signature := c.GetHeader("Authorization")
	if signature == "" {
		reqresp.ResponseMarshal(c, -1, "need sign", &resp)
		return
	}

	if ok, err := checkAuth(req.Common.AppId, appSecret, signature, &req); !ok {
		reqresp.ResponseMarshal(c, -2, err.Error(), &resp)
		return
	}

	ipInfo, err := ip2region.Inst().MemorySearch(req.ClientIP)
	if err != nil {
		reqresp.ResponseMarshal(c, -3, err.Error(), &resp)
		return
	}
	resp.ClientIP = req.ClientIP
	resp.Country = ipInfo.Country
	resp.City = ipInfo.City
	resp.Region = ipInfo.Region
	resp.Province = ipInfo.Province
	resp.ISP = ipInfo.ISP
	resp.CityId = ipInfo.CityId
	reqresp.ResponseMarshal(c, errors.OK.Code, errors.OK.Msg, nil)
}