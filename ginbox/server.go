package ginbox

import (
	"context"
	"github.com/joyous-x/enceladus/common/xlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	default_server_name = "default"
	default_server_addr = ":8000"
)

// GinServer GinServer 对象
type GinServer struct {
	engine   *gin.Engine
	httpsvr  *http.Server
	signal   chan int
	name     string
	addr     string
	certFile string
	keyFile  string
}

// GinServerBox GinServerBox 对象
type GinServerBox struct {
	middlewares []gin.HandlerFunc
	servers     map[string]*GinServer
}

// MakeGinServerBox 快速的根据配置生成gin server
// 初始化时, 可以制定一些中间件作为 gin sever 的通用中间件
// 返回 GinServerBox 对象指针
func MakeGinServerBox(middleware ...gin.HandlerFunc) *GinServerBox {
	tmp := &GinServerBox{
		middlewares: middleware,
	}
	return tmp
}

// Server 获取name指定的配置相关联的gin.Engine
// note: when len(names) == 0, return default_server
func (this *GinServerBox) Server(names ...string) *gin.Engine {
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

	return this.servers[name].engine
}

// NewServer 根据参数生成一个gin.Engine
// 一般用作内部封装使用，用于生成 GinServerBox 所包含的 server 对象
func (this *GinServerBox) NewServer(name, addr, certFile, keyFile string, debug bool, middleware ...gin.HandlerFunc) (*gin.Engine, error) {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	if len(name) == 0 {
		name = default_server_name
	}

	if len(addr) == 0 {
		addr = default_server_addr
	}

	if strings.Index(addr, ":") < 0 {
		return nil, fmt.Errorf("invalid addr: should be (host):port")
	}

	if nil == this.servers {
		this.servers = make(map[string]*GinServer)
	} else {
		if v, ok := this.servers[name]; ok {
			return v.engine, fmt.Errorf("already exists")
		}
	}

	router := gin.New()
	router.Use(this.middlewares...)
	router.Use(middleware...)
	this.servers[name] = &GinServer{
		engine:   router,
		name:     name,
		addr:     addr,
		certFile: certFile,
		keyFile:  keyFile,
	}

	return router, nil
}

// Handle 根据参数处理server的route相关
func (this *GinServerBox) Handle(name, method, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	if len(name) == 0 {
		name = default_server_name
	}

	v, ok := this.servers[name]
	if !ok {
		return nil
	}

	return v.engine.Handle(method, relativePath, handlers...)
}

// RunServers 启动并运行各个server
// 此函数会阻塞，直到各个server退出
func (this *GinServerBox) RunServers() error {
	if this.servers == nil || len(this.servers) < 1 {
		return fmt.Errorf("no server")
	}

	defer this.StopServers()
	for k, _ := range this.servers {
		err := this.runSingleServer(k)
		if err != nil {
			return err
		}
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	return nil
}

// StopServers 停止各个server
// 通常会等待各个server完全响应终止信号
func (this *GinServerBox) StopServers() error {
	if this.servers == nil || len(this.servers) < 1 {
		return nil
	}
	for k, _ := range this.servers {
		this.stopSingleServer(k)
	}

	for k, v := range this.servers {
		<-v.signal
		xlog.Info("http server has stopped: %v", k)
	}
	this.servers = nil

	return nil
}

func (this *GinServerBox) stopSingleServer(name string) (err error) {
	v, ok := this.servers[name]
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = v.httpsvr.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		xlog.Warn("http server Shutdown: %v", err)
	}
	return
}

func (this *GinServerBox) runSingleServer(name string) (err error) {
	v, ok := this.servers[name]
	if !ok {
		return
	}

	httpSvr := &http.Server{
		Addr:    v.addr, // fmt.Sprintf(":%d", port)
		Handler: v.engine,
	}
	sigint := make(chan int, 1)

	this.servers[name].httpsvr = httpSvr
	this.servers[name].signal = sigint
	go func() {
		if len(v.certFile) > 0 && len(v.keyFile) > 0 {
			xlog.Debug("GinServerBox server: ready to ListenAndServeTLS: %s(%s) certFile=%v keyFile=%v", name, v.addr, v.certFile, v.keyFile)
			if err := httpSvr.ListenAndServeTLS(v.certFile, v.keyFile); err != http.ErrServerClosed {
				xlog.Warn("GinServerBox server: ListenAndServeTLS: %v", err)
			}
		} else {
			xlog.Debug("GinServerBox server: ready to ListenAndServe: %s(%s) certFile=%v keyFile=%v", name, v.addr, v.certFile, v.keyFile)
			if err := httpSvr.ListenAndServe(); err != http.ErrServerClosed {
				// Error starting or closing listener:
				xlog.Warn("GinServerBox server: ListenAndServe: %v", err)
			}
		}
		sigint <- 1
	}()

	xlog.Info("GinServerBox server: %s(%s) running", name, v.addr)
	return
}
