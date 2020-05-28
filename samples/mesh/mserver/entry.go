package mserver

import (
	"context"
	"enceladus/common/xlog"
	mproto "enceladus/project/mesh_demo/protoc"
	"enceladus/wgrpc"
	"enceladus/wgrpc/protoc"
	wproto "enceladus/wgrpc/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
)

var entry = &entryServer{}

// FirstEntry ...
func FirstEntry(port string) error {
	xlog.Debug("===> FirstEntry ready")
	model.init()

	svcRegistFunc := func(s *wgrpc.WServer) error {
		mproto.RegisterEntryServer(s.Server, entry)
		return nil
	}
	server := wgrpc.NewDispatchServer(entrySvcName, port, svcRegistFunc)
	if nil == server {
		err := fmt.Errorf("serverfunc.NewDispatchServer error")
		xlog.Error("FirstEntry error: %v", err)
		return err
	}
	server.Start()

	xlog.Debug("s----- FirstEntry end")
	return nil
}

////////////////////////////////////// grpc server: dispatcher /////////////////////////////////////////

func init() {
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Entry_Insert), EntryInsertHandler)
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Entry_Sum), EntrySumHandler)
	wgrpc.RegistDispatchHander(int32(mproto.C2S_Cmds_Entry_Max), EntryMaxHandler)
}

// EntryInsertHandler ...
func EntryInsertHandler(ctx context.Context, comm *wproto.RouteHeader, reqData []byte) (*wproto.DispatchResp, error) {
	inreq := &mproto.InsertReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := insertEntryInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

// EntrySumHandler ...
func EntrySumHandler(ctx context.Context, comm *wproto.RouteHeader, reqData []byte) (*wproto.DispatchResp, error) {
	inreq := &mproto.SumReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := sumEntryInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

// EntryMaxHandler ...
func EntryMaxHandler(ctx context.Context, comm *wproto.RouteHeader, reqData []byte) (*wproto.DispatchResp, error) {
	inreq := &mproto.MaxReq{}
	resp, _ := NewDispatchResp(ctx, comm)

	err := proto.Unmarshal(reqData, inreq)
	if err != nil {
		return resp, err
	}

	inresp, err := maxEntryInner(ctx, inreq)
	if err != nil {
		return resp, err
	}

	return fillDispatchResp(resp, inresp)
}

func fillDispatchResp(resp *wproto.DispatchResp, msg proto.Message) (*wproto.DispatchResp, error) {
	respPbData, err := proto.Marshal(msg)
	if err != nil {
		return resp, err
	}
	resp.RetCode = 0
	resp.Pbdata = respPbData
	return resp, nil
}

func insertEntryInner(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	return entry.Insert(ctx, req)
}

func sumEntryInner(ctx context.Context, req *mproto.SumReq) (*mproto.SumResp, error) {
	return entry.Sum(ctx, req)
}

func maxEntryInner(ctx context.Context, req *mproto.MaxReq) (*mproto.MaxResp, error) {
	return entry.Max(ctx, req)
}

///////////////////////////////////////// grpc server ////////////////////////////////////////////

// entryServer implementation of grpc interface EntryServer
type entryServer struct {
}

// Insert func of EntryServer
func (t *entryServer) Insert(ctx context.Context, req *mproto.InsertReq) (*mproto.InsertResp, error) {
	resp := &mproto.InsertResp{}
	reqPBdata, _ := proto.Marshal(req)

	reqDispatch := &protoc.DispatchReq{
		Header: &protoc.RouteHeader{
			ReqId: fmt.Sprintf("%v", time.Now().Unix()),
			Uid:   "test_uid",
			Appid: "test_appid",
			Cmd:   int32(mproto.C2S_Cmds_Middle_InsertDecorator),
		},
		Pbdata: reqPBdata,
	}

	respDispatch, err := wgrpc.CallDispatch(ctx, middleSvcName, "", reqDispatch)
	if err != nil {
		xlog.Error("entryServer.Insert.CallDispatch req=%+v err=%v", *req, err)
		return resp, err
	}

	err = proto.Unmarshal(respDispatch.Pbdata, resp)
	if err != nil {
		return resp, err
	}

	xlog.Debug("entryServer.Insert req=%+v resp=%+v respDispatch=%+v", *req, *resp, *respDispatch)
	return resp, nil
}

// Sum func of entryServer
func (t *entryServer) Sum(ctx context.Context, req *mproto.SumReq) (*mproto.SumResp, error) {
	resp := &mproto.SumResp{}
	reqPBdata, _ := proto.Marshal(req)

	reqDispatch := &protoc.DispatchReq{
		Header: &protoc.RouteHeader{
			ReqId: fmt.Sprintf("%v", time.Now().Unix()),
			Uid:   "test_uid",
			Appid: "test_appid",
			Cmd:   int32(mproto.C2S_Cmds_Model_Sum),
		},
		Pbdata: reqPBdata,
	}

	respDispatch, err := wgrpc.CallDispatch(ctx, modelSvcName, "", reqDispatch)
	if err != nil {
		xlog.Error("entryServer.Sum.CallDispatch req=%+v err=%v", *req, err)
		return resp, err
	}

	err = proto.Unmarshal(respDispatch.Pbdata, resp)
	if err != nil {
		return resp, err
	}

	xlog.Debug("entryServer.Sum req=%+v resp=%+v", *req, *resp)
	return resp, nil
}

// Max func of entryServer
func (t *entryServer) Max(ctx context.Context, req *mproto.MaxReq) (*mproto.MaxResp, error) {
	resp := &mproto.MaxResp{}
	svcGrpcName := "com.mesh_demo.protoc.Model"
	err := wgrpc.CallRawRe(ctx, modelSvcName, "", svcGrpcName, "Max", req, resp)
	if err != nil {
		xlog.Error("entryServer.Max.CallRawRe req=%+v err=%v", *req, err)
		return resp, err
	}
	xlog.Debug("entryServer.Max req=%+v resp=%+v", *req, *resp)
	return resp, nil
}
