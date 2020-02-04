package router

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func version(c *gin.Context) {
	datas := map[string]string{
		"data":   "hello world",
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}
	if "GET" == c.Request.Method {
		datas["args"] = c.Query("args")
	} else if "POST" == c.Request.Method {
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, datas)
}
