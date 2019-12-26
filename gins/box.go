package gins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"sync"
)

var default_ginbox *GinBox

func init () {
	default_ginbox = func(middleware ...gin.HandlerFunc) *GinBox {
		return &GinBox{
			middlewares: middleware,
		}
	}(DefaultMiddlewares...)
}

// Default 默认的全局GinServerBox对象
func DefaultBox() *GinBox {
	return default_ginbox
}

// GinBox GinBox 对象
type GinBox struct {
	middlewares []gin.HandlerFunc
	servers     map[string]*GinServer
}

// InitDefault 初始化默认的全局GinServerBox对象
func (this *GinBox) Init(configs []*ServerConfig, middleware ...gin.HandlerFunc) error {
	this.servers = make(map[string]*GinServer, 0)
	this.middlewares = append(this.middlewares, middleware...)
	for i, v := range configs {
		if len(v.Name) == 0 {
			v.Name = default_server_name
		}
		if v.Port == 0 {
			v.Port = default_server_port
		}
		s, err := this.newServer(v, this.middlewares...)
		if err != nil {
			xlog.Panic("InitServers pos:%v, name:%v, port:%v, err:%v", i, v.Name, v.Port, err)
		} else {
			this.servers[v.Name] = s
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
func (this *GinBox) newServer(conf *ServerConfig, middleware ...gin.HandlerFunc) (*GinServer, error) {
	server := NewGinServer()
	return server, server.Init(conf, middleware...)
}

// Handle 根据参数处理server的route相关
func (this *GinBox) Handle(name, method, relativePath string, handlers ...gin.HandlerFunc) error {
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
		xlog.Error("GinBox Run error: invalid servers")
		return fmt.Errorf("no servers")
	}

	var waitGroup sync.WaitGroup
	for _, s := range this.servers {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
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
