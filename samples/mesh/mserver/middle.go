package mserver

import (
	"context"
	"enceladus/common/xlog"
	mproto "enceladus/project/mesh_demo/protoc"
	"enceladus/wgrpc"
	wproto "enceladus/wgrpc/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var middle = &middleServer{}

// MiddleEntry ...
func MiddleEntry(port string) error {
	xlog.Debug("===> MiddleEntry ready")
	model.init()

	svcRegistFunc := func(s *wgrpc.WServer) error {
		mproto.RegisterMiddleServer(s.Server, middle)
		return nil
	}
	server := wgrpc.NewDispatchServer(middleSvcName, port, svcRegistFunc)
	if nil == server {
		err := fmt.Errorf("serverfunc.NewDispatchServer error")
		xlog.Error("MiddleEntry error: %v", err)
		return err
	}
	server.Start()

	xlog.Debug("s----- MiddleEntry end")
	return nil
}

////////////////////////////////////// grpc server: dispatcher /////////////////////////////////////////

func init() {
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Middle_InsertDecorator), MiddleInsertDeHandler)
}

// MiddleInsertDeHandler ...
func MiddleInsertDeHandler(ctx context.Context, comm *wproto.RouteHeader, reqData []byte) (*wproto.DispatchResp, error) {
	inreq := &mproto.InsertReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := insertDecoratorInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

func insertDecoratorInner(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	return middle.InsertDecorator(ctx, req)
}

///////////////////////////////////////// grpc server ////////////////////////////////////////////

// middleServer implementation of grpc interface MiddleServer
type middleServer struct {
}

// middleServer InsertDecorator func of middleServer
func (t *middleServer) InsertDecorator(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	caller := func(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
		modelClient := mproto.NewModelClient(conn)
		return modelClient.Insert(ctx, req)
	}
	resp, err := wgrpc.CallRaw(ctx, modelSvcName, "", caller)
	if err != nil {
		return nil, err
	}
	xlog.Debug("middleServer.InsertDecorator req=%+v resp=%+v", *req, *(resp.(*mproto.InsertResp)))
	return resp.(*mproto.InsertResp), nil
}
