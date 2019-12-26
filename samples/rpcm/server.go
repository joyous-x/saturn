package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/rpcm"
	"github.com/joyous-x/saturn/rpcm/base"
	"net/http"
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

func main() {
	xlog.Debug("rpcm.HTTP sample ===> start ")
	conf := &base.ServiceConfig{
		Protocal: "http",
		Name:     "",
		Port:     8001,
	}
	iserver, err := rpcm.NewService(conf)
	if err != nil {
		xlog.Error("new service err:%v config:%+v", err, *conf)
		return
	}
	iserver.Route("POST", "/v1", version)
	iserver.Run()
	xlog.Debug("rpcm.HTTP sample ===> end ")
}
