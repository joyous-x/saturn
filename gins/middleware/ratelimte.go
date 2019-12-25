package middleware

import (
	"sync/atomic"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/gin-gonic/gin"
)

type RateLimitConcurrency struct {
	maxConcurrency int64
	current        int64
}

// RateLimitOnConcurrency 限制最大并发数
// 估算最大并发数公式：
// 		单机峰值OPS * 接口平均延时,
// 例如：
//		服务峰值ops = 8000个请求每秒, 接口平均延时为 25ms
//	那么,
// 		最大并发数 = 8000 * 0.025 = 200
func RateLimitOnConcurrency(maxConcurrency int64, logOnly bool) gin.HandlerFunc {
	rateLimiter := &RateLimitConcurrency{
		maxConcurrency: maxConcurrency,
		current: 0,
	}
	return func(c *gin.Context) {
		atomic.AddInt64(&rateLimiter.current, 1)
		defer atomic.AddInt64(&rateLimiter.current, -1)

		current := atomic.LoadInt64(&rateLimiter.current)
		if current > rateLimiter.maxConcurrency {
			if logOnly {
				xlog.Warn("%v is busy", c.Request.URL.Path)
			} else {
				c.JSON(200, errors.ErrRateLimit)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}