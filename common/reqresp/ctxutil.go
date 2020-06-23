package reqresp

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CtxKeyClientIP       = "ctx.client.ip"
	CtxKeyUuid           = "ctx.client.uuid"
	CtxKeyAppID          = "ctx.client.appid"
	CtxKeyVersion        = "ctx.client.version"
	CtxKeyDeviceID       = "ctx.client.deviceid"
	CtxKeyRequestID      = "ctx.client.rid"
	CtxKeyRequestTime    = "ctx.client.rtimestamp"
	CtxKeyEchoToken      = "ctx.client.echo_token"
	CtxKeyMetricName     = "ctx.g.metric_name"
	CtxKeyRequestData    = "ctx.request.data"
	CtxKeyResponseData   = "ctx.response.data"
	CtxKeyResponseObject = "ctx.response.object"
	CtxKeySessionToken   = "ctx.session.token"
)

// F 打印出带上下文的日志
func F(ctx context.Context, format string, args ...interface{}) string {
	ip, _ := ctx.Value(CtxKeyClientIP).(string)
	uuid, _ := ctx.Value(CtxKeyUuid).(string)
	rid, _ := ctx.Value(CtxKeyRequestID).(string)
	appID, _ := ctx.Value(CtxKeyAppID).(string)
	requestTime, _ := ctx.Value(CtxKeyRequestTime).(time.Time)
	version, _ := ctx.Value(CtxKeyVersion).(string)
	metricName, _ := ctx.Value(CtxKeyMetricName).(string)

	t := requestTime.Format("20060102")

	fullArgs := make([]interface{}, 0, 7+len(args))
	fullArgs = append(fullArgs, metricName, rid, uuid, appID, ip, version, t)
	fullArgs = append(fullArgs, args...)

	return fmt.Sprintf("[m=%s rid=%s uuid=%s app=%s ip=%s ver=%s t=%s]"+format, fullArgs...)
}

// MakeCtx 返回标准格式的 context
func MakeCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, CtxKeyClientIP, c.ClientIP())
	ctx = context.WithValue(ctx, CtxKeyUuid, c.GetString(CtxKeyUuid))
	ctx = context.WithValue(ctx, CtxKeyAppID, c.GetString(CtxKeyAppID))
	ctx = context.WithValue(ctx, CtxKeyVersion, c.GetString(CtxKeyVersion))
	ctx = context.WithValue(ctx, CtxKeyDeviceID, c.GetString(CtxKeyDeviceID))
	ctx = context.WithValue(ctx, CtxKeyRequestID, c.GetString(CtxKeyRequestID))
	ctx = context.WithValue(ctx, CtxKeyRequestTime, c.GetTime(CtxKeyRequestTime))
	if d, exist := c.Get(CtxKeyRequestData); exist {
		if rawData, ok := d.([]byte); ok {
			ctx = context.WithValue(ctx, CtxKeyRequestData, rawData)
		}
	}
	return ctx
}

// CtxSetRaw ...
func CtxSetRaw(c *gin.Context, key string, data []byte) {
	if data != nil {
		c.Set(key, data)
	}
}

// CtxGetRaw ...
func CtxGetRaw(c context.Context, key string) []byte {
	if v := c.Value(key); v != nil {
		if rawData, ok := v.([]byte); ok {
			return rawData
		}
	}
	return nil
}

// CtxGetStr ...
func CtxGetStr(c context.Context, key string) string {
	if v := c.Value(key); v != nil {
		if data, ok := v.(string); ok {
			return data
		}
	}
	return ""
}
