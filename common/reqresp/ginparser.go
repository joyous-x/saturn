package reqresp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"

	"io/ioutil"
	"net/http"
	"time"
	// "github.com/mitchellh/mapstructure"
)

// RequestUnmarshal 解析入参
func RequestUnmarshal(c *gin.Context, data IRequest) (context.Context, error) {
	return requestUnmarshal(c, data)
}

func requestUnmarshal(c *gin.Context, data IRequest) (context.Context, error) {
	ctx := context.Background()

	reqRawPacket := CtxGetRaw(c, CtxKeyRequestData)
	if nil == reqRawPacket {
		reqRawPacket, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return ctx, err
		}
		CtxSetRaw(c, CtxKeyRequestData, reqRawPacket)
	}

	err := json.Unmarshal(reqRawPacket, &data)
	if err != nil {
		xlog.Error("Unmarshal error raw=%v %v\n", string(reqRawPacket), err)
		return ctx, err
	}

	appid, uid, sessToken := "", "", ""
	if ireq, ok := data.(IRequest); ok {
		ireqCom := ireq.GetCommon()
		appid, uid = ireqCom.AppID, ireqCom.Uid
		c.Set(CtxKeyDeviceID, ireqCom.DeviceID)
		c.Set(CtxKeyEchoToken, ireqCom.EchoToken)
		c.Set(CtxKeyAppID, appid)
		c.Set(CtxKeyUuid, uid)
		c.Set(CtxKeyRequestID, fmt.Sprintf("%x-%x", time.Now().Unix(), utils.NewSafeRand(time.Now().Unix()).Int31()))
		c.Set(CtxKeySessionToken, sessToken)
		ctx = MakeCtx(c)
	}
	return ctx, err
}

// ResponseMarshal 序列化应答数据 s.BaseError
func ResponseMarshal(c *gin.Context, err error, data IResponse) {
	var berr = errors.OK
	if err != nil {
		tmp, ok := err.(errors.BaseError)
		if !ok {
			berr = errors.NewError(errors.ErrServerError.Code, err.Error())
		} else {
			berr = tmp
		}
	}
	responseMarshal(c, berr.Code, berr.Msg, data, http.StatusOK)
}

func responseMarshal(c *gin.Context, status int, message string, data IResponse, httpcode int) {
	if data == nil {
		data = &RespCommon{}
	}
	iresp, ok := data.(IResponse)
	if !ok {
		panic("invalid response data: not IResponse")
	}
	common := iresp.GetCommon()
	common.Ret = status
	common.Msg = message
	common.RequestID = c.MustGet(CtxKeyRequestID).(string)
	common.EchoToken = c.MustGet(CtxKeyEchoToken).(string)
	common.RetryMS = 1000
	common.Timestamp = time.Now().Unix()

	c.Set("response_time", common.Timestamp)
	c.Set("ret", common.Ret)
	c.Set("msg", common.Msg)
	responseBody, _ := json.Marshal(data)
	if token, exist := c.Get(CtxKeySessionToken); exist {
		signature := utils.MakeHMac(token.(string), responseBody)
		c.Writer.Header().Set("Authorization", signature)
	}
	c.Data(httpcode, "application/json", responseBody)
	return
}
