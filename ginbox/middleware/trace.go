package middleware

import (
	"github.com/joyous-x/enceladus/common/xlog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
)

func newUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// Trace 拦截器
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ctx.client.rid", newUUID())
		c.Set("ctx.client.rtimestamp", time.Now())
		xlog.Info("gin: %v, %v, %v,  start", c.Request.Method, c.Request.URL.Path, c.Request.URL)
		defer func() {
			xlog.Info("gin: %v, %v, %v,  end", c.Request.Method, c.Request.URL.Path, c.Request.URL)
		}()
		c.Next()
	}
}
