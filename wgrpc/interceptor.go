package wgrpc

import (
	"context"
	_ "github.com/grpc-ecosystem/go-grpc-middleware" // ...
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/govern/tracing"
	"google.golang.org/grpc"
	"time"
)

// GenRecoverInterceptor ...
func GenRecoverInterceptor(handler grpc_recovery.RecoveryHandlerFunc) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	grpcRecoverOpt := []grpc_recovery.Option{grpc_recovery.WithRecoveryHandler(handler)}
	unaryRecover := grpc_recovery.UnaryServerInterceptor(grpcRecoverOpt...)
	streamRecover := grpc_recovery.StreamServerInterceptor(grpcRecoverOpt...)
	return unaryRecover, streamRecover
}

// UnaryServerInterceptorAccessLog ...
func UnaryServerInterceptorAccessLog() grpc.UnaryServerInterceptor {
	unaryServerInterceptorAccessLog := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return resp, err
		}
		escaped := time.Now().Sub(start).Nanoseconds() / 1000
		if err != nil {
			xlog.Error("UnaryServerInterceptorAccessLog method=%v escaped_μs=%v error: %v", info.FullMethod, escaped, err)
		} else {
			xlog.Debug("UnaryServerInterceptorAccessLog method=%v escaped_μs=%v", info.FullMethod, escaped)
		}
		return resp, err
	}
	return unaryServerInterceptorAccessLog
}

// StreamServerInterceptorAccessLog ...
func StreamServerInterceptorAccessLog() grpc.StreamServerInterceptor {
	streamServerInterceptorAccessLog := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		escaped := time.Now().Sub(start).Nanoseconds() / 1000
		if err != nil {
			xlog.Error("StreamServerInterceptorAccessLog method=%v escaped_μs=%v error: %v", info.FullMethod, escaped, err)
		} else {
			xlog.Debug("StreamServerInterceptorAccessLog method=%v escaped_μs=%v", info.FullMethod, escaped)
		}
		return err
	}
	return streamServerInterceptorAccessLog
}

// UnaryServerInterceptorJaeger ...
func UnaryServerInterceptorJaeger() grpc.UnaryServerInterceptor {
	return tracing.JaegerUnaryServerInterceptor(tracing.JaegerGlobalTracer())
}

// UnaryClientInterceptorJaeger ...
func UnaryClientInterceptorJaeger() grpc.UnaryClientInterceptor {
	return tracing.JaegerUnaryClientInterceptor(tracing.JaegerGlobalTracer())
}
