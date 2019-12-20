package middleware

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		whiteList := [...]string{".baidu.com", ".google.com"}
		req_url_str := c.GetHeader("Referer")
		req_origin := c.GetHeader("Origin")
		{
			req_host := req_origin
			is_valid := false
			for _, s := range whiteList {
				if strings.Contains(req_host, s) {
					is_valid = true
					break
				}
			}
			if is_valid {
				xlog.Debug("[CORS]req_url='%s', req_host='%s' in whiteList:%v, host:%s,%s, %v", req_url_str, req_host, is_valid, c.Request.Host, c.Request.URL, c.Request)
				header := c.Writer.Header()
				header.Set("Access-Control-Allow-Origin", req_host)
				header.Set("Access-Control-Allow-Credentials", "true")
				header.Set("Access-Control-Allow-Methods", "GET, POST")
				header.Set("Access-Control-Allow-Headers", "Keep-Alive,User-Agent,Cache-Control,Content-Type,Authorization")
				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
					return
				}
				if c.Request.Method == "HEAD" {
					c.AbortWithStatus(200)
					return
				}
			}
		}
		c.Next()
	}
}

var CorsEx = cors.New(cors.Config{
	AllowMethods:    []string{"GET,POST,PUT,PATCH,DELETE,OPTIONS"},
	AllowHeaders:    []string{"Keep-Alive,User-Agent,Cache-Control,Content-Type,Authorization"},
	AllowAllOrigins: true,
})
