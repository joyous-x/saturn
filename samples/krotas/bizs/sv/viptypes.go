package sv

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
)

// VipTypesReq ...
type VipTypesReq struct {
	reqresp.ReqCommon
}

// VipTypesResp ...
type VipTypesResp struct {
	reqresp.RespCommon
	Items []VipTypeItem `json:"vip_types"`
}

// VipTypeItem ...
type VipTypeItem struct {
	Level       int    `json:"level"`
	PriceAli    int    `json:"price_ali"`
	PriceWechat int    `json:"price_wechat"`
	Desc        string `json:"desc"`
}

// VipTypes ...
func VipTypes(c *gin.Context) {
	req := VipTypesReq{}
	resp := VipTypesResp{}
	_, err := reqresp.RequestUnmarshal(c, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, &resp)
		return
	}

	resp.Items = append(resp.Items, VipTypeItem{1, 250, 250, "test_1"})
	resp.Items = append(resp.Items, VipTypeItem{2, 750, 750, "test_2"})
	resp.Items = append(resp.Items, VipTypeItem{3, 1500, 1500, "test_3"})

	reqresp.ResponseMarshal(c, errors.OK, &resp)
}
