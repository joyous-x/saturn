package mserver

import (
	"container/heap"
	"context"
	"enceladus/common/xlog"
	mproto "enceladus/project/mesh_demo/protoc"
	""enceladus/common/utils""
	"enceladus/wgrpc"
	"enceladus/wgrpc/protoc"
	"fmt"
	proto "github.com/golang/protobuf/proto"
)

var model = &modelServer{}

// ModelEntry ...
func ModelEntry(port string) error {
	xlog.Debug("===> ModelEntry ready")
	model.init()

	svcRegistFunc := func(s *wgrpc.WServer) error {
		mproto.RegisterModelServer(s.Server, model)
		return nil
	}
	server := wgrpc.NewDispatchServer(modelSvcName, port, svcRegistFunc)
	if nil == server {
		err := fmt.Errorf("serverfunc.NewDispatchServer error")
		xlog.Error("ModelEntry error: %v", err)
		return err
	}
	server.Start()

	xlog.Debug("s----- ModelEntry end")
	return nil
}

////////////////////////////////////// grpc server: dispatcher /////////////////////////////////////////

func init() {
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Model_Sum), ModelSumHandler)
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Model_Max), ModelMaxHandler)
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Model_Len), ModelLenHandler)
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Model_Insert), ModelInsertHandler)
}

// ModelInsertHandler ...
func ModelInsertHandler(ctx context.Context, comm *protoc.RouteHeader, reqData []byte) (*protoc.DispatchResp, error) {
	inreq := &mproto.InsertReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := insertModelInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

// ModelSumHandler dispatcher's handler
func ModelSumHandler(ctx context.Context, comm *protoc.RouteHeader, reqData []byte) (*protoc.DispatchResp, error) {
	inreq := &mproto.SumReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := sumModelInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

// ModelMaxHandler dispatcher's handler
func ModelMaxHandler(ctx context.Context, comm *protoc.RouteHeader, reqData []byte) (*protoc.DispatchResp, error) {
	inreq := &mproto.MaxReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := maxModelInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

// ModelLenHandler dispatcher's handler
func ModelLenHandler(ctx context.Context, comm *protoc.RouteHeader, reqData []byte) (*protoc.DispatchResp, error) {
	inreq := &mproto.LenReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := lenModelInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

func insertModelInner(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	return model.Insert(ctx, req)
}
func sumModelInner(ctx context.Context, req *mproto.SumReq) (*mproto.SumResp, error) {
	return model.Sum(ctx, req)
}
func maxModelInner(ctx context.Context, req *mproto.MaxReq) (*mproto.MaxResp, error) {
	return model.Max(ctx, req)
}
func lenModelInner(ctx context.Context, req *mproto.LenReq) (*mproto.LenResp, error) {
	return model.Len(ctx, req)
}

///////////////////////////////////////// grpc server ////////////////////////////////////////////

// modelServer implementation of grpc interface ModelServer
type modelServer struct {
	intHeap *utils.IntHeap
}

func (m *modelServer) init() error {
	m.intHeap = &utils.IntHeap{}
	heap.Init(m.intHeap)
	return nil
}

// Insert func of modelServer
func (m *modelServer) Insert(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	resp := &mproto.InsertResp{}
	heap.Push(m.intHeap, req.Val)
	resp.Len = int32(m.intHeap.Len())
	resp.Data = m.intHeap.Data()
	xlog.Debug("modelServer.Insert req=%+v resp=%+v", *req, *resp)
	return resp, nil
}

// Sum func of modelServer
func (m *modelServer) Sum(ctx context.Context, req *mproto.SumReq) (*mproto.SumResp, error) {
	resp := &mproto.SumResp{Sum: 0}
	for _, v := range m.intHeap.Data() {
		resp.Sum += int32(v)
	}
	xlog.Debug("modelServer.Sum req=%+v resp=%+v", *req, *resp)
	return resp, nil
}

// Max func of modelServer
func (m *modelServer) Max(ctx context.Context, req *mproto.MaxReq) (*mproto.MaxResp, error) {
	resp := &mproto.MaxResp{}
	dat := heap.Pop(m.intHeap)
	if dat == nil {
		return resp, fmt.Errorf("pop error")
	}
	resp.Max = int32(dat.(int32))
	xlog.Debug("modelServer.Max req=%+v resp=%+v", *req, *resp)
	return resp, nil
}

// Len func of modelServer
func (m *modelServer) Len(ctx context.Context, req *mproto.LenReq) (*mproto.LenResp, error) {
	resp := &mproto.LenResp{}
	resp.Len = int32(m.intHeap.Len())
	xlog.Debug("modelServer.Len req=%+v resp=%+v", *req, *resp)
	return resp, nil
}
