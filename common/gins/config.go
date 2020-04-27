package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/gins/middleware"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/rpcm/base"
)

type ServerConfig = base.ServiceConfig

const (
	defaultServerName = "default"
	defaultServerPort = 8000
)

// DefaultMiddlewares ...
var DefaultMiddlewares = []gin.HandlerFunc{middleware.Trace(), middleware.Recovery(RecoveryHandler), middleware.Cors()}

// RecoveryHandler recovery中间件的默认处理函数
func RecoveryHandler(c *gin.Context) {
	xlog.Warn("recovery handler: %v, %v", c.Request.Method, c.Request.URL)
}

// TranGinHandlerFunc2Interface ...
func TranGinHandlerFunc2Interface(datas []gin.HandlerFunc) []interface{} {
	result := make([]interface{}, len(datas))
	for i, d := range datas {
		result[i] = d
	}
	return result
}
