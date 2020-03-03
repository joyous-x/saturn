package common

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/common/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	// "github.com/mitchellh/mapstructure"
)

const (
	sessionToken = "session_token"
	defaultToken = "ffffffffffffffffffffff"
)

// RequestBody 统一请求头
type RequestBody struct {
	UUID      string          `json:"uuid"`
	AppName   string          `json:"appname"`
	Timestamp int64           `json:"timestamp"`
	Version   string          `json:"version"`
	DeviceID  string          `json:"device_id"`
	EchoToken string          `json:"echo_token"`
	Data      json.RawMessage `json:"data"` // 业务数据
}

// ResponseBody 统一响应头
type ResponseBody struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id"`
	RetryMS   int32       `json:"retry_ms"`
	EchoToken string      `json:"echo_token"`
	Data      interface{} `json:"data"` // 业务数据
}

// FnAuthUserInfo 获取用户信息的回调
type FnAuthUserInfo func(appname, uuid string) (string, error)

// RequestUnmarshal 解析入参
func RequestUnmarshal(c *gin.Context, fn FnAuthUserInfo, data interface{}) (ctx context.Context, appName, uuid string, err error) {
	return requestUnmarshal(c, fn, data, true)
}

// RequestUnmarshalNoAuth 解析入参
func RequestUnmarshalNoAuth(c *gin.Context, fn FnAuthUserInfo, data interface{}) (ctx context.Context, appName, uuid string, err error) {
	return requestUnmarshal(c, fn, data, false)
}

func requestUnmarshal(c *gin.Context, fnUserInfo FnAuthUserInfo, data interface{}, isNeedCheck bool) (ctx context.Context, appName, uuid string, err error) {
	requestRawPacket, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	requestBody := RequestBody{}
	err = json.Unmarshal(requestRawPacket, &requestBody)
	if err != nil {
		xlog.Error("Unmarshal error raw=%v %v\n", string(requestRawPacket), err)
		return
	}

	uuid = requestBody.UUID
	appName = requestBody.AppName
	c.Set(DeviceID, requestBody.DeviceID)
	c.Set(EchoToken, requestBody.EchoToken)
	c.Set(Version, requestBody.Version)
	c.Set(AppName, appName)
	c.Set(Uuid, uuid)
	ctx = MakeCtx(c)

	if data != nil {
		err = json.Unmarshal(requestBody.Data, data)
		if err != nil {
			xlog.Error("requestUnmarshal error: %v", err)
			return
		}
	}

	var token string
	token, err = fnUserInfo(appName, uuid)
	if err != nil {
		xlog.Error("RetrieveToken error: %v", err)
		return
	}

	if isNeedCheck && token == "" {
		err = fmt.Errorf("Authorization check fail")
		return
	}

	//signature
	signSha1, signDevice := parseAuthorization(c.GetHeader("Authorization"))
	signDeviceServer := utils.Fingerprint(string(requestRawPacket))

	xlog.Info("signature app:%s, uuid:%s, client:%s, server:%s, body:%s", appName, uuid, signDevice, signDeviceServer, string(requestRawPacket))

	serverSign := authSignContent(token, requestRawPacket)
	if isNeedCheck && serverSign != signSha1 {
		xlog.Error("signature client:%s, server:%s, token:%s, body:%s", signSha1, serverSign, token, string(requestRawPacket))
		err = fmt.Errorf("Authorization check fail")
		return
	}

	c.Set(sessionToken, token)
	return
}

// MakeAuthSign 签名算法
func MakeAuthSign(token string, body []byte) string {
	return authSignContent(token, body)
}

func authSignContent(token string, body []byte) (signature string) {
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write(body)
	signature = hex.EncodeToString(mac.Sum(nil))
	return
}

// ResponseMarshal 序列化应答数据
func ResponseMarshal(c *gin.Context, status int, message string, data interface{}) {
	responseMarshalInner(c, status, message, data, http.StatusOK)
}

// ResponseMarshalWithCode 序列化应答数据
func ResponseMarshalWithCode(c *gin.Context, status int, message string, data interface{}, httpcode int) {
	responseMarshalInner(c, status, message, data, httpcode)
}

// ResponseMarshal 压缩传参
func responseMarshalInner(c *gin.Context, status int, message string, data interface{}, httpcode int) {
	responseBody := ResponseBody{}
	responseBody.Status = status
	responseBody.Message = message
	responseBody.EchoToken = c.MustGet(EchoToken).(string)
	responseBody.RequestID = c.MustGet(RequestId).(string)
	responseBody.RetryMS = 1000
	responseBody.Timestamp = time.Now().Unix()
	if data == nil {
		responseBody.Data = map[string]string{}
	} else {
		responseBody.Data = data
	}
	token, exists := c.Get(sessionToken)

	c.Set("response_time", responseBody.Timestamp)
	c.Set("status", responseBody.Status)
	c.Set("message", responseBody.Message)

	responseBodyData, _ := json.Marshal(responseBody)
	if exists {
		signature := authSignContent(token.(string), responseBodyData)
		c.Writer.Header().Set("Authorization", signature)
	}
	c.Data(httpcode, "application/json", responseBodyData)
	return
}

func parseAuthorization(allAuth string) (string, string) {
	//Http Header -> Authorization: WOW sign="ewrqew2", csign="45436ggsdg"
	re := regexp.MustCompile(`WOW\s*sign="(\w+)",\s*csign="(\w+)"`)
	params := re.FindStringSubmatch(allAuth)
	if len(params) == 3 {
		return params[1], params[2]
	}

	re = regexp.MustCompile(`WOW\s*(\w+)`)
	params = re.FindStringSubmatch(allAuth)
	if len(params) == 2 {
		return params[1], ""
	}

	return "", ""
}
