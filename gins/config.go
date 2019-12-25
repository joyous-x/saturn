package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/gins/middleware"
)

type ServerConfig struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	Debug    bool   `yaml:"debug"`
}

const (
	default_server_name = "default"
	default_server_port = 8000
)

var default_middlewares = []gin.HandlerFunc{middleware.Trace(), middleware.Recovery(RecoveryHandler), middleware.Cors()}

// RecoveryHandler recovery中间件的默认处理函数
func RecoveryHandler(c *gin.Context) {
	xlog.Warn("recovery handler: %v, %v", c.Request.Method, c.Request.URL)
}
