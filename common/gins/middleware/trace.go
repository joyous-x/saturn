package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joyous-x/saturn/common/xlog"
	"strings"
	"time"
)

func newUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// Trace 拦截器
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceBegin := time.Now()
		c.Set("ctx.client.rid", newUUID())
		c.Set("ctx.client.rtimestamp", traceBegin)
		xlog.Info("gin: %v, %v, %v  start", c.Request.Method, c.Request.URL.Path, c.Request.URL)
		defer func() {
			xlog.Info("gin: %v, %v, %v, spend(μs)=%v  end", c.Request.Method, c.Request.URL.Path, c.Request.URL, (time.Now().UnixNano()-traceBegin.UnixNano())/1000000)
		}()
		c.Next()
	}
}
