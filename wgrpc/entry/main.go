package main

import (
	"context"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/wgrpc"
	eprotoc "github.com/joyous-x/saturn/wgrpc/entry/protoc"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

var (
	server = "127.0.0.1"
	port   = 8086
)

// CalcServerImpl implention of CalcServer
type CalcServerImpl struct{}

// Sum sum function for CalcServer
func (s *CalcServerImpl) Sum(ctx context.Context, req *eprotoc.SumReq) (*eprotoc.SumResp, error) {
	return sumInner(ctx, req)
}

func sumInner(ctx context.Context, req *eprotoc.SumReq) (*eprotoc.SumResp, error) {
	resp := &eprotoc.SumResp{}
	resp.S = req.A / req.B
	xlog.Debug("------ CalcServer.Sum: %v + %v => %v", req.A, req.B, resp.S)
	return resp, nil
}

func clientFuncInsecure() {
	ctx := context.Background()
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", server, port), grpc.WithInsecure())
	if err != nil {
		xlog.Error("c----- clientfunc.Dial: err %v", err)
	}
	defer conn.Close()
	cc := eprotoc.NewCalcClient(conn)

	for i := 0; i < 10000; i++ {
		req := &eprotoc.SumReq{A: int32(i), B: 1}
		resp, err := cc.Sum(ctx, req)
		if err != nil {
			xlog.Error("c----- clientfunc.Sum: err %v", err)
		} else {
			xlog.Debug("c----- clientfunc.Sum: resp=%v", resp.S)
		}
		time.Sleep(100 * time.Millisecond)
	}

	xlog.Debug("c----- clientfunc end")
}

func serverFunc() {
	svcInfo := &wgrpc.ServerInfo{
		ConnTimeoutSec: 5,
		UniqID:         "wgrpc_entry_main",
		Name:           "wgrpc_entry_main",
	}
	server, err := wgrpc.NewWServer(svcInfo, "", "")
	if nil != err {
		xlog.Error("s----- serverfunc error: %v", err)
		return
	}
	eprotoc.RegisterCalcServer(server.Server, &CalcServerImpl{})
	server.Start()
	xlog.Debug("s----- serverfunc end")
}

func main() {
	t := flag.String("t", "c", "c or s, cd or sd: client or server, dispatch client or dispatch server")
	p := flag.Int("p", 8086, "port for client")
	flag.Parse()
	xlog.Debug("------ grpc mode: %v", *t)

	port = *p
	if *t == "c" {
		clientFuncInsecure()
	} else if *t == "s" {
		serverFunc()
	} else if *t == "cd" {
		clientDispatchFuncInsecure()
	} else if *t == "sd" {
		serverDispatchFunc()
	} else {
		flag.PrintDefaults()
	}
	return
}
