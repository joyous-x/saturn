package common

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	ClientIp       = "ctx.client.ip"
	Uuid           = "ctx.client.uuid"
	AppName        = "ctx.client.appname"
	Version        = "ctx.client.version"
	DeviceID       = "ctx.client.deviceid"
	RequestId      = "ctx.client.rid"
	RequestTime    = "ctx.client.rtimestamp"
	EchoToken      = "ctx.client.echo_token"
	MetricName     = "ctx.g.metric_name"
	RequestData    = "ctx.request.data"
	ResponseData   = "ctx.response.data"
	ResponseObject = "ctx.response.object"
)

// F 打印出带上下文的日志
func F(ctx context.Context, format string, args ...interface{}) string {
	ip, _ := ctx.Value(ClientIp).(string)
	uuid, _ := ctx.Value(Uuid).(string)
	rid, _ := ctx.Value(RequestId).(string)
	appName, _ := ctx.Value(AppName).(string)
	requestTime, _ := ctx.Value(RequestTime).(time.Time)
	version, _ := ctx.Value(Version).(string)
	metricName, _ := ctx.Value(MetricName).(string)

	t := requestTime.Format("20060102")

	fullArgs := make([]interface{}, 0, 7+len(args))
	fullArgs = append(fullArgs, metricName, rid, uuid, appName, ip, version, t)
	fullArgs = append(fullArgs, args...)

	return fmt.Sprintf("[m=%s rid=%s uuid=%s app=%s ip=%s ver=%s t=%s]"+format, fullArgs...)
}

// MakeCtx 返回标准格式的 context
func MakeCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, RequestId, c.GetString(RequestId))
	ctx = context.WithValue(ctx, RequestTime, c.GetTime(RequestTime))
	ctx = context.WithValue(ctx, Version, c.GetString(Version))
	ctx = context.WithValue(ctx, DeviceID, c.GetString(DeviceID))
	ctx = context.WithValue(ctx, ClientIp, c.ClientIP())
	ctx = context.WithValue(ctx, Uuid, c.GetString(Uuid))
	ctx = context.WithValue(ctx, AppName, c.GetString(AppName))
	return ctx
}
