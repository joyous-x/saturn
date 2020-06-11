package bizs

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/reqresp"
)

// AliPayReq ...
type AliPayReq struct {
	reqresp.ReqCommon
}

// AliPayResp ...
type AliPayResp struct {
	reqresp.RespCommon
}

// AliPay ali pay
func AliPay(c *gin.Context) {
	//> TODO
}
