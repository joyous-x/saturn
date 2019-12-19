package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func version(c *gin.Context) {
	datas := map[string]string{"data": "hello world"}
	if "GET" == c.Request.Method {
		datas["args"] = c.Query("args")
	} else if "POST" == c.Request.Method {

	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	panic("xxxxxx")
	c.JSON(http.StatusOK, datas)
}
