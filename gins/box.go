package ginbox

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"sync"
)

var default_ginbox = MakeGinBox(default_middlewares...)

// Default 默认的全局GinServerBox对象
func DefaultBox() *GinBox {
	return default_ginbox
}

// GinBox GinBox 对象
type GinBox struct {
	middlewares []gin.HandlerFunc
	servers     map[string]*GinServer
}

// MakeGinBox 快速的根据配置生成gin server
// 初始化时, 可以制定一些中间件作为 gin sever 的通用中间件
// 返回 GinBox 对象指针
func MakeGinBox(middleware ...gin.HandlerFunc) *GinBox {
	tmp := &GinBox{
		middlewares: middleware,
	}
	return tmp
}

// InitDefault 初始化默认的全局GinServerBox对象
func (this *GinBox) Init(configs []*ServerConfig, middleware ...gin.HandlerFunc) error {
	this.middlewares = append(this.middlewares, middleware...)
	for i, v := range configs {
		_, err := this.newServer(v, this.middlewares...)
		if err != nil {
			xlog.Panic("InitServers pos:%v, name:%v, port:%v, err:%v", i, v.Name, v.Port, err)
		}
	}
	return nil
}

// Server 获取name指定的配置相关联的gin.Engine
// note: when len(names) == 0, return server named default_server_name
func (this *GinBox) Server(names ...string) *GinServer {
	if nil == this.servers {
		return nil
	}
	name := default_server_name
	if len(names) > 0 {
		name = names[0]
	}
	if _, ok := this.servers[name]; !ok {
		return nil
	}
	return this.servers[name]
}

// newServer 根据参数生成一个gin.Engine
func (this *GinBox) newServer(conf *ServerConfig, middleware ...gin.HandlerFunc) (*gin.Engine, error) {
	server := NewGinServer()
	return server.Init(conf, middleware...)
}

// Handle 根据参数处理server的route相关
func (this *GinBox) Handle(name, method, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	if len(name) == 0 {
		name = default_server_name
	}
	v, ok := this.servers[name]
	if !ok {
		return nil
	}
	return v.Handle(method, relativePath, handlers...)
}

// Run 启动并运行各个server
// 此函数会阻塞，直到各个server退出
func (this *GinBox) Run() error {
	if this.servers == nil || len(this.servers) < 1 {
		return fmt.Errorf("no servers")
	}

	var waitGroup sync.WaitGroup
	for _, s := range this.servers {
		go func() {
			defer waitGroup.Done()
			waitGroup.Add(1)
			s.Run()
		}()
	}

	waitGroup.Wait()
	return nil
}

// Stop 停止各个server
// 通常会等待各个server完全响应终止信号
func (this *GinBox) Stop() error {
	if this.servers == nil || len(this.servers) < 1 {
		return nil
	}
	for k, s := range this.servers {
		s.Stop()
		xlog.Info("http server has stopped: %v", k)
	}
	this.servers = nil
	return nil
}
