package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"runtime/debug"
)

func Recovery(deal func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				xlog.Error("[Recovery]path='%s' err='%s' stack='%s'", c.Request.URL, err, stack)
				if nil != deal {
					deal(c)
				}
			}
		}()
		c.Next()
	}
}
