package wgrpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// reference:
//          https://github.com/grpc/grpc/blob/master/doc/health-checking.md

func registHealthSvc(server *grpc.Server) error {
	hs := health.NewServer()
	healthpb.RegisterHealthServer(server, hs)
	return nil
}
