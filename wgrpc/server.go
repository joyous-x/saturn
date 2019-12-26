package wgrpc

import (
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/govern/wconsul"
	"net"
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/grpc-ecosystem/go-grpc-middleware/recovery" // ...
	"github.com/joyous-x/saturn/govern/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// ServerInfo ...
type ServerInfo struct {
	UniqID           string
	Name             string
	Host             string
	Port             int
	ConnTimeoutSec   int
	Tags             []string
	AddrConsulCenter string
	AddrJaegerAgent  string
}

var defaultGrpcRecoveryHandlerFunc = func(p interface{}) (err error) {
	xlog.Error("recover from panic. error:  %v", p)
	xlog.Error("===> panic stack:  %v", string(debug.Stack()))
	return fmt.Errorf("%v", p)
}

// WServer  warpper of grpc server
type WServer struct {
	*grpc.Server
	srv      *ServerInfo
	keyFile  string
	certFile string
}

func makeGrpcInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	unaryRecoverInterceptor, streamRecoverInterceptor := GenRecoverInterceptor(defaultGrpcRecoveryHandlerFunc)
	unaryInterceptors := []grpc.UnaryServerInterceptor{unaryRecoverInterceptor}
	streamInterceptors := []grpc.StreamServerInterceptor{streamRecoverInterceptor}

	//> append other interceptors here
	unaryInterceptors = append(unaryInterceptors, UnaryServerInterceptorAccessLog())
	streamInterceptors = append(streamInterceptors, StreamServerInterceptorAccessLog())

	unaryInterceptors = append(unaryInterceptors, UnaryServerInterceptorJaeger())

	return unaryInterceptors, streamInterceptors
}

// NewWServer  new instance of WServer
func NewWServer(srv *ServerInfo, certFile, keyFile string, opt ...grpc.ServerOption) (*WServer, error) {
	wserver := &WServer{srv: srv}
	if len(srv.UniqID) < 10 {
		srv.UniqID = fmt.Sprintf("%v_%v:%v_%v", srv.Name, srv.Host, srv.Port, time.Now().Unix())
	}

	if "" != certFile && "" != keyFile {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			xlog.Error("Failed to generate credentials %v", err)
			return nil, err
		}
		opt = append(opt, grpc.Creds(creds))
		wserver.certFile, wserver.keyFile = certFile, keyFile
	}

	_, _, err := tracing.NewJaegerTracer(srv.Name, srv.AddrJaegerAgent)
	if err != nil {
		xlog.Warn("fail new jaeger tracer, err=%v", err)
	}

	opt = append(opt, grpc.ConnectionTimeout(time.Second*time.Duration(srv.ConnTimeoutSec)))
	unaryInterceptors, streamInterceptors := makeGrpcInterceptors()
	opt = append(opt, grpc_middleware.WithUnaryServerChain(unaryInterceptors...), grpc_middleware.WithStreamServerChain(streamInterceptors...))
	wserver.Server = grpc.NewServer(opt...)

	xlog.Debug("===> new wserver ready: UniqID=%v port=%v connTimeoutSec=%v creds=%v key=%v", srv.UniqID, srv.Port, srv.ConnTimeoutSec, certFile, keyFile)
	return wserver, nil
}

// Start start to serve
func (s *WServer) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", s.srv.Port))
	if err != nil {
		return err
	}

	useTLS := func() bool {
		if len(s.certFile) > 0 && len(s.keyFile) > 0 {
			return true
		}
		return false
	}()
	svcInfo := &wconsul.SvcInfo{
		ID:      s.srv.UniqID,
		Name:    s.srv.Name,
		Host:    s.srv.Host,
		Port:    s.srv.Port,
		Tags:    append([]string{}, s.srv.Tags...),
		Checker: wconsul.NewGrpcSvcCheck(s.srv.Host, s.srv.Port, "5s", "3s", useTLS),
	}
	if err := wconsul.SvcRegistration(svcInfo); err != nil {
		// s.Server.Stop()
		xlog.Warn("===> WServer.Start.SvcRegistration uniqID=%v err=%v", s.srv.UniqID, err)
	}
	defer func() {
		xlog.Info("===> WServer.Start.SvcDeregistration uniqID=%v ready", svcInfo.ID)
		wconsul.SvcDeregistration(svcInfo.ID)
	}()

	registHealthSvc(s.Server)
	reflection.Register(s.Server)

	return s.Serve(l)
}
