package gins

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
)

var defaultGinbox *GinBox

// IGinBoxRouter interface for gins's router
type IGinBoxRouter interface {
	HTTPRouter(ginbox *GinBox) error
}

func init() {
	defaultGinbox = func(middleware ...gin.HandlerFunc) *GinBox {
		return &GinBox{
			middlewares: middleware,
		}
	}(DefaultMiddlewares...)
}

// DefaultBox 默认的全局GinServerBox对象
func DefaultBox() *GinBox {
	return defaultGinbox
}

// GinBox GinBox 对象
type GinBox struct {
	middlewares []gin.HandlerFunc
	servers     map[string]*GinServer
}

// Init 初始化默认的全局GinServerBox对象
func (g *GinBox) Init(configs []*ServerConfig, middleware ...gin.HandlerFunc) error {
	g.servers = make(map[string]*GinServer, 0)
	g.middlewares = append(g.middlewares, middleware...)
	for i, v := range configs {
		if len(v.Name) == 0 {
			v.Name = defaultServerName
		}
		if v.Port == 0 {
			v.Port = defaultServerPort
		}
		s, err := g.newServer(v, g.middlewares...)
		if err != nil {
			xlog.Panic("InitServers pos:%v, name:%v, port:%v, err:%v", i, v.Name, v.Port, err)
		} else {
			g.servers[v.Name] = s
		}
	}

	return nil
}

// Server 获取name指定的配置相关联的gin.Engine
// note: when len(names) == 0, return server named defaultServerName
func (g *GinBox) Server(names ...string) *GinServer {
	if nil == g.servers {
		return nil
	}
	name := defaultServerName
	if len(names) > 0 {
		name = names[0]
	}
	if _, ok := g.servers[name]; !ok {
		return nil
	}
	return g.servers[name]
}

// newServer 根据参数生成一个gin.Engine
func (g *GinBox) newServer(conf *ServerConfig, middleware ...gin.HandlerFunc) (*GinServer, error) {
	server := NewGinServer()
	return server, server.Init(conf, middleware...)
}

// Handle 根据参数处理server的route相关
func (g *GinBox) Handle(name, method, relativePath string, handlers ...gin.HandlerFunc) error {
	if len(name) == 0 {
		name = defaultServerName
	}
	v, ok := g.servers[name]
	if !ok {
		return nil
	}
	return v.Handle(method, relativePath, handlers...)
}

// HTTPRouter regist a http router
func (g *GinBox) HTTPRouter(irouter IGinBoxRouter) error {
	return irouter.HTTPRouter(g)
}

// Run 启动并运行各个server
// 此函数会阻塞，直到各个server退出
func (g *GinBox) Run() error {
	if g.servers == nil || len(g.servers) < 1 {
		xlog.Error("GinBox Run error: invalid servers")
		return fmt.Errorf("no servers")
	}

	var waitGroup sync.WaitGroup
	for i := range g.servers {
		waitGroup.Add(1)
		index := i
		go func() {
			defer waitGroup.Done()
			g.servers[index].Run()
		}()
	}

	waitGroup.Wait()
	return nil
}

// Stop 停止各个server
// 通常会等待各个server完全响应终止信号
func (g *GinBox) Stop() error {
	if g.servers == nil || len(g.servers) < 1 {
		return nil
	}
	for k, s := range g.servers {
		s.Stop()
		xlog.Info("http server has stopped: %v", k)
	}
	g.servers = nil
	return nil
}
