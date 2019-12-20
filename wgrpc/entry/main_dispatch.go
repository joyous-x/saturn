package main

import (
	"context"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/wgrpc"
	eprotoc "github.com/joyous-x/saturn/wgrpc/entry/protoc"
	"github.com/joyous-x/saturn/wgrpc/protoc"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"time"
)

const (
	cmdCalcAdd = 1
)

func init() {
	wgrpc.RegistDispatchHander(cmdCalcAdd, CalcSumHandler)
}

// CalcSumHandler ...
func CalcSumHandler(ctx context.Context, comm *protoc.RouteHeader, reqData []byte) (*protoc.DispatchResp, error) {
	sumreq := &eprotoc.SumReq{}
	resp := &protoc.DispatchResp{
		Header: &protoc.RouteHeader{
			ReqId: comm.ReqId,
			Uid:   comm.Uid,
			Appid: comm.Appid,
			Cmd:   comm.Cmd,
		},
		Cmd:     comm.Cmd + 1,
		RetCode: -1,
	}

	err := proto.Unmarshal(reqData, sumreq)
	if err != nil {
		return resp, err
	}

	sumresp, err := sumInner(ctx, sumreq)
	if err != nil {
		return resp, err
	}

	respPbData, err := proto.Marshal(sumresp)
	if err != nil {
		return resp, err
	}

	resp.RetCode = 0
	resp.Pbdata = respPbData
	return resp, nil
}

func clientDispatchFuncInsecure() {
	ctx := context.Background()
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", server, port), grpc.WithInsecure())
	if err != nil {
		xlog.Error("c----- clientDispatchFuncInsecure.Dial: err %v", err)
		return
	}
	defer conn.Close()
	cc := protoc.NewDispatchSvcClient(conn)

	sumpbdata, _ := proto.Marshal(&eprotoc.SumReq{A: 3, B: 0})
	sumresp := &eprotoc.SumResp{}

	req := &protoc.DispatchReq{
		Header: &protoc.RouteHeader{
			ReqId: fmt.Sprintf("%v", time.Now().Unix()),
			Uid:   "test_uid",
			Appid: "test_appid",
			Cmd:   cmdCalcAdd,
		},
		Pbdata: sumpbdata,
	}
	resp, err := cc.Dispatch(ctx, req)
	if err != nil {
		xlog.Error("c----- clientDispatchFuncInsecure.Dispatch: err %v", err)
		return
	}
	err = proto.Unmarshal(resp.Pbdata, sumresp)
	if err != nil {
		xlog.Error("c----- clientDispatchFuncInsecure.Unmarshal: err %v", err)
		return
	}

	xlog.Debug("c----- clientDispatchFuncInsecure.Dispatch: header=%v resp_data=%v", resp.Header, sumresp)
	time.Sleep(100 * time.Millisecond)
	xlog.Debug("c----- clientDispatchFuncInsecure end")
}

func serverDispatchFunc() {
	svcRegistFunc := func(s *wgrpc.WServer) error {
		eprotoc.RegisterCalcServer(s.Server, &CalcServerImpl{})
		return nil
	}
	server := wgrpc.NewDispatchServer("", "", svcRegistFunc)
	if nil == server {
		xlog.Error("s----- serverfunc.NewDispatchServer error")
		return
	}
	server.Start()
	xlog.Debug("s----- serverfunc end")
}
