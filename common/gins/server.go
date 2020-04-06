package gins

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// GinServer GinServer 对象
type GinServer struct {
	ServerConfig
	engine      *gin.Engine
	httpsvr     *http.Server
	signStop    chan int
	middlewares []gin.HandlerFunc
}

// NewGinServer 快速的根据配置生成gin server
// 初始化时, 可以制定一些中间件作为 gin sever 的通用中间件
// 返回 GinServer 对象指针
func NewGinServer(middleware ...gin.HandlerFunc) *GinServer {
	tmp := &GinServer{
		middlewares: middleware,
	}
	return tmp
}

// Server 获取指定配置相关联的gin.Engine
func (this *GinServer) Engine() *gin.Engine {
	return this.engine
}

// InitServer 根据参数生成一个gin.Engine
func (this *GinServer) Init(conf *ServerConfig, middleware ...gin.HandlerFunc) error {
	name, port, certFile, keyFile, debug := conf.Name, conf.Port, conf.CertFile, conf.KeyFile, conf.Debug

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	if len(name) == 0 {
		name = default_server_name
	}
	if port == 0 {
		port = default_server_port
	}

	router := gin.New()
	router.Use(this.middlewares...)
	router.Use(middleware...)
	this.engine = router
	this.Name = name
	this.Port = port
	this.CertFile = certFile
	this.KeyFile = keyFile
	this.middlewares = append(this.middlewares, middleware...)

	return nil
}

// Handle 根据参数处理server的route相关
func (this *GinServer) Handle(method, relativePath string, handlers ...gin.HandlerFunc) error {
	var err error
	if this.engine != nil {
		iRoutes := this.engine.Handle(method, relativePath, handlers...)
		if iRoutes != nil {
			return nil
		} else {
			err = fmt.Errorf("route handler error(%v, %v)", method, relativePath)
		}
	} else {
		err = fmt.Errorf("server not ready: (%v, %v)", method, relativePath)
	}
	return err
}

// Handle 根据参数处理server的route相关
func (this *GinServer) Route(method, relativePath string, routes ...interface{}) error {
	var err error
	if this.engine != nil {
		handlers := make([]gin.HandlerFunc, len(routes))
		for i, r := range routes {
			if f, ok := r.(gin.HandlerFunc); ok {
				handlers[i] = f
			} else if f, ok := r.(func(*gin.Context)); ok {
				handlers[i] = gin.HandlerFunc(f)
			} else {
				xlog.Panic("route handler error(invalid handler) %v %v", method, relativePath)
			}
		}
		iRoutes := this.engine.Handle(method, relativePath, handlers...)
		if iRoutes != nil {
			return nil
		} else {
			err = fmt.Errorf("route handler error(%v, %v)", method, relativePath)
		}
	} else {
		err = fmt.Errorf("server not ready: (%v, %v)", method, relativePath)
	}
	return err
}

// Run 启动并运行server
// 此函数会阻塞，直到server退出
func (this *GinServer) Run() error {
	err := this.runServer()
	if err != nil {
		return err
	}

	defer this.Stop()
	signOSInter := make(chan os.Signal, 1)
	signal.Notify(signOSInter, os.Interrupt)
WAIT:
	select {
	case <-signOSInter:
	case <-this.signStop:
	default:
		goto WAIT
	}

	return nil
}

func (this *GinServer) runServer() error {
	httpSvr := &http.Server{
		Addr:    fmt.Sprintf(":%v", this.Port),
		Handler: this.engine,
	}
	this.httpsvr = httpSvr
	this.signStop = make(chan int, 1)

	name, addr := this.Name, httpSvr.Addr
	certFile, keyFile := this.CertFile, this.KeyFile

	go func() {
		if len(certFile) > 0 && len(keyFile) > 0 {
			xlog.Debug("GinServer: ready to ListenAndServeTLS: %s(%s) certFile=%v keyFile=%v", name, addr, certFile, keyFile)
			if err := httpSvr.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
				xlog.Warn("GinServer:%s(%s) ListenAndServeTLS: %v", name, addr, err)
			}
		} else {
			xlog.Debug("GinServer ready to ListenAndServe: %s(%s) certFile=%v keyFile=%v", name, addr, certFile, keyFile)
			if err := httpSvr.ListenAndServe(); err != http.ErrServerClosed {
				// Error starting or closing listener:
				xlog.Warn("GinServer:%s(%s) ListenAndServe: %v", name, addr, err)
			}
		}
		xlog.Debug("GinServer:%s(%s) ListenAndServe stopped", name, addr)
		this.signStop <- 1
	}()

	xlog.Info("GinServer:%s(%s) running", name, addr)
	return nil
}

// Stop 停止server, 通常会等待server完全响应终止信号
func (this *GinServer) Stop() error {
	if this.httpsvr == nil {
		return nil
	}
	err := this.stopServer(5 * time.Second)
	if err != nil {
		return err
	}
	<-this.signStop
	this.engine = nil
	return nil
}

func (this *GinServer) stopServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := this.httpsvr.Shutdown(ctx)
	if err != nil {
		// Error from closing listeners, or context timeout:
		xlog.Warn("http server Shutdown: %v", err)
	}
	return err
}
