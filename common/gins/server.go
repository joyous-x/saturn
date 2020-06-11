package gins

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/xlog"
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

// Engine 获取指定配置相关联的gin.Engine
func (g *GinServer) Engine() *gin.Engine {
	return g.engine
}

// Init 根据参数生成一个gin.Engine
func (g *GinServer) Init(conf *ServerConfig, middleware ...gin.HandlerFunc) error {
	name, port, certFile, keyFile, debug := conf.Name, conf.Port, conf.CertFile, conf.KeyFile, conf.Debug

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	if len(name) == 0 {
		name = defaultServerName
	}
	if port == 0 {
		port = defaultServerPort
	}

	router := gin.New()
	router.Use(g.middlewares...)
	router.Use(middleware...)
	g.ReadTimeoutMs = conf.ReadTimeoutMs
	g.WriteTimeoutMs = conf.WriteTimeoutMs
	g.engine = router
	g.Name = name
	g.Port = port
	g.CertFile = certFile
	g.KeyFile = keyFile
	g.middlewares = append(g.middlewares, middleware...)

	return nil
}

// Handle 根据参数处理server的route相关
func (g *GinServer) Handle(method, relativePath string, handlers ...gin.HandlerFunc) error {
	var err error
	if g.engine != nil {
		iRoutes := g.engine.Handle(method, relativePath, handlers...)
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

// Route 根据参数处理server的route相关
func (g *GinServer) Route(method, relativePath string, routes ...interface{}) error {
	var err error
	if g.engine != nil {
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
		iRoutes := g.engine.Handle(method, relativePath, handlers...)
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
func (g *GinServer) Run() error {
	err := g.runServer()
	if err != nil {
		return err
	}

	defer g.Stop()
	signOSInter := make(chan os.Signal, 1)
	signal.Notify(signOSInter, os.Interrupt)
WAIT:
	select {
	case <-signOSInter:
	case <-g.signStop:
	default:
		goto WAIT
	}

	return nil
}

func (g *GinServer) runServer() error {
	httpSvr := &http.Server{
		Addr:         fmt.Sprintf(":%v", g.Port),
		Handler:      g.engine,
		ReadTimeout:  g.ReadTimeoutMs * time.Millisecond,
		WriteTimeout: g.WriteTimeoutMs * time.Millisecond,
	}
	g.httpsvr = httpSvr
	g.signStop = make(chan int, 1)

	name, addr := g.Name, httpSvr.Addr
	certFile, keyFile := g.CertFile, g.KeyFile

	go func() {
		if len(certFile) > 0 && len(keyFile) > 0 {
			xlog.Debug("GinServer: ready to ListenAndServeTLS: %s(%s) certFile=%v keyFile=%v", name, addr, certFile, keyFile)
			if err := httpSvr.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				xlog.Warn("GinServer:%s(%s) ListenAndServeTLS: %v", name, addr, err)
			}
		} else {
			xlog.Debug("GinServer ready to ListenAndServe: %s(%s) certFile=%v keyFile=%v", name, addr, certFile, keyFile)
			if err := httpSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// Error starting or closing listener:
				xlog.Warn("GinServer:%s(%s) ListenAndServe: %v", name, addr, err)
			}
		}
		xlog.Debug("GinServer:%s(%s) ListenAndServe stopped", name, addr)
		g.signStop <- 1
	}()

	xlog.Info("GinServer:%s(%s) running", name, addr)
	return nil
}

// Stop 停止server, 通常会等待server完全响应终止信号
func (g *GinServer) Stop() error {
	if g.httpsvr == nil {
		return nil
	}
	err := g.stopServer(5 * time.Second)
	if err != nil {
		return err
	}
	<-g.signStop
	g.engine = nil
	return nil
}

func (g *GinServer) stopServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := g.httpsvr.Shutdown(ctx)
	if err != nil {
		// Error from closing listeners, or context timeout:
		xlog.Warn("http server Shutdown: %v", err)
	}
	return err
}
