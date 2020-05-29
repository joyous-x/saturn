package reqresp

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
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

const (
	defaultToken = "ffffffffffffffffffffff"
)

// FnAuthUser 认证用户信息, 认证通过时返回session token的函数定义
type FnAuthUser func(appid, uid, authorization string, dataBody []byte) (string, error)

// RequestUnmarshal 解析入参
func RequestUnmarshal(c *gin.Context, fn FnAuthUser, data IRequest) (ctx context.Context, err error) {
	return requestUnmarshal(c, fn, data)
}

func requestUnmarshal(c *gin.Context, fnAuthUser FnAuthUser, data IRequest) (ctx context.Context, err error) {
	requestRawPacket, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(requestRawPacket, &data)
	if err != nil {
		xlog.Error("Unmarshal error raw=%v %v\n", string(requestRawPacket), err)
		return
	}

	appid, uid, sessToken := "", "", ""
	if ireq, ok := data.(IRequest); ok {
		ireqCom := ireq.GetCommon()
		appid, uid = ireqCom.AppId, ireqCom.Uid
		c.Set(DeviceID, ireqCom.DeviceID)
		c.Set(EchoToken, ireqCom.EchoToken)
		c.Set(AppId, appid)
		c.Set(Uuid, uid)
		c.Set(RequestId, fmt.Sprintf("%x-%x", time.Now().Unix(), utils.NewSafeRand(time.Now().Unix()).Int31()))
		ctx = MakeCtx(c)
	}

	if nil != fnAuthUser {
		sessToken, err = fnAuthUser(appid, uid, c.GetHeader("Authorization"), requestRawPacket)
		if err != nil {
			return
		}
	}

	c.Set(SessionToken, sessToken)
	return
}

// ResponseMarshal 序列化应答数据
func ResponseMarshal(c *gin.Context, err errors.BaseError, data IResponse) {
	responseMarshal(c, err.Code, err.Msg, data, http.StatusOK)
}

func responseMarshal(c *gin.Context, status int, message string, data IResponse, httpcode int) {
	common := func() *RespCommonData {
		if data == nil {
			return &RespCommonData{}
		}
		iresp, ok := data.(IResponse)
		if !ok {
			panic("invalid response data: not IResponse")
		}
		return iresp.GetCommon()
	}()

	common.Ret = status
	common.Msg = message
	common.RequestId = c.MustGet(RequestId).(string)
	common.EchoToken = c.MustGet(EchoToken).(string)
	common.RetryMS = 1000
	common.Timestamp = time.Now().Unix()

	c.Set("response_time", common.Timestamp)
	c.Set("ret", common.Ret)
	c.Set("msg", common.Msg)

	token, exists := c.Get(SessionToken)
	responseBody, _ := json.Marshal(data)
	if exists {
		signature := authSignContent(token.(string), responseBody)
		c.Writer.Header().Set("Authorization", signature)
	}
	c.Data(httpcode, "application/json", responseBody)
	return
}

////////////////////////////////////

func authSignContent(token string, body []byte) (signature string) {
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write(body)
	signature = hex.EncodeToString(mac.Sum(nil))
	return
}

func authUserSample(appid, uid, authorization string, dataBody []byte) (string, error) {
	// get user token
	token := ""

	serverSign := authSignContent(token, dataBody)
	if serverSign != authorization {
		xlog.Error("signature client:%s, server:%s, token:%s, body:%s", authorization, serverSign, token, string(dataBody))
		return "", fmt.Errorf("Authorization check fail")
	}

	return token, nil
}
