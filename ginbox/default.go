package ginbox

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/ginbox/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Name        string
	Port        int
	CertFile    string
	KeyFile     string
	Debug       bool
	Middlewares []gin.HandlerFunc
}

var default_middlewares = []gin.HandlerFunc{middleware.Trace(), middleware.Recovery(RecoveryHandler), middleware.Cors()}
var default_ginserverbox = MakeGinServerBox(default_middlewares...)

// RecoveryHandler recovery中间件的默认处理函数
func RecoveryHandler(c *gin.Context) {
	xlog.Warn("recovery handler: %v, %v", c.Request.Method, c.Request.URL)
}

// Default 默认的全局GinServerBox对象
func Default() *GinServerBox {
	return default_ginserverbox
}

// DefaultServer 默认的全局GinServerBox对象的默认server
func DefaultServer() *gin.Engine {
	return default_ginserverbox.Server()
}

// InitDefault 初始化默认的全局GinServerBox对象
func InitDefault(configs []*ServerConfig) (*GinServerBox, error) {
	var err error
	for i, v := range configs {
		middlewares := func() []gin.HandlerFunc {
			if nil == v.Middlewares {
				return make([]gin.HandlerFunc, 0)
			}
			return v.Middlewares
		}()
		_, err = default_ginserverbox.NewServer(v.Name, fmt.Sprintf(":%v", v.Port), v.CertFile, v.KeyFile, v.Debug, middlewares...)
		if err != nil {
			xlog.Panic("InitServers pos:%v, name:%v, port:%v, err:%v", i, v.Name, v.Port, err)
			break
		}
	}
	return default_ginserverbox, err
}
