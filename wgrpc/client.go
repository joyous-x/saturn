package wgrpc

import (
	"google.golang.org/grpc"
)

type WGrpcClient struct {
	*grpc.ClientConn
}
