package wgrpc

import (
	"context"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/govern/wconsul"
	"github.com/joyous-x/saturn/wgrpc/protoc"
	"fmt"
	"google.golang.org/grpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	cretFile    = ""
	cretKey     = ""
	svcCenter   = "127.0.0.1:8500"
	jaegerAgent = "127.0.0.1:6831"
)

// InitModules ...
func InitModules(svcCenterAddr, jaegerAgentAddr string) error {
	if len(svcCenterAddr) > 0 {
		svcCenter = svcCenterAddr
	}
	if len(jaegerAgentAddr) > 0 {
		jaegerAgent = jaegerAgentAddr
	}

	if consulClient := wconsul.InitDefaultClient(svcCenter); len(svcCenter) > 0 && consulClient == nil {
		err := fmt.Errorf("invalid consul client: not ready")
		xlog.Error("NewWServer error: %v", err)
		return err
	}

	return nil
}

// SvcRegistFunc in this function, we should regist services to grpc server
type SvcRegistFunc func(s *WServer) error

// DispatchHandler handler for DispatchSvc
type DispatchHandler func(context.Context, *protoc.RouteHeader, []byte) (*protoc.DispatchResp, error)

var once sync.Once
var dispatchSvc *dispatchSvcImpl

func defaultDispatchSvc() *dispatchSvcImpl {
	once.Do(func() {
		dispatchSvc = &dispatchSvcImpl{}
		dispatchSvc.handlers = make(map[int32]DispatchHandler, 0)
	})
	return dispatchSvc
}

// RegistDispatchHander regist handlers for dispatcher
func RegistDispatchHander(cmd int32, handler DispatchHandler) error {
	return defaultDispatchSvc().RegistHandler(cmd, handler)
}

// NewDispatchServer generate a grpc server which has a service named DispatchSvc
func NewDispatchServer(name, port string, register SvcRegistFunc, opt ...grpc.ServerOption) *WServer {
	svcInfo := &ServerInfo{
		ConnTimeoutSec:   5,
		Name:             name,
		AddrConsulCenter: svcCenter,
	}
	if len(port) < 1 || strings.Index(port, ":") < 0 {
		freePort, err := utils.GetLocalFreePort()
		if err != nil {
			xlog.Error("NewDispatchServer.GetLocalFreePort error: %v", err)
			return nil
		}
		svcInfo.Port = freePort
	} else {
		hp := strings.Split(port, ":")
		if len(hp) != 2 {
			xlog.Error("NewDispatchServer invalid port: %v", port)
			return nil
		}
		if nport, err := strconv.Atoi(hp[1]); err == nil {
			svcInfo.Host, svcInfo.Port = hp[0], nport
		} else {
			xlog.Error("NewDispatchServer port=%v err: %v", port, err)
			return nil
		}
	}

	//> TODO: get cret and key
	certFile, cretKey := "", ""

	wserver, err := NewWServer(svcInfo, certFile, cretKey)
	if err != nil {
		xlog.Error("NewDispatchServer.NewWServer error: %v", err)
		return nil
	}

	protoc.RegisterDispatchSvcServer(wserver.Server, defaultDispatchSvc())
	if err = register(wserver); err != nil {
		xlog.Error("NewDispatchServer.SvcRegistFunc error: %v", err)
		wserver.Server.Stop()
		return nil
	}

	return wserver
}

type dispatchSvcImpl struct {
	handlers map[int32]DispatchHandler
}

func (s *dispatchSvcImpl) RegistHandler(cmd int32, handler DispatchHandler) error {
	if _, ok := s.handlers[cmd]; ok {
		xlog.Warn("RegistHandler rewrite handler of cmd=%v", cmd)
	}
	s.handlers[cmd] = handler
	return nil
}

func (s *dispatchSvcImpl) Dispatch(ctx context.Context, req *protoc.DispatchReq) (*protoc.DispatchResp, error) {
	start := time.Now()
	found := false

	defer func() {
		xlog.Info("Dispatch(cmd_found=%v) header=%v escaped_μs=%v", found, req.Header, time.Since(start).Nanoseconds()/1000)
	}()
	if handler, ok := s.handlers[req.Header.Cmd]; ok {
		found = true
		return handler(ctx, req.Header, req.Pbdata)
	}

	return nil, fmt.Errorf("do not find handler for cmd=%v", req.Header.Cmd)
}
