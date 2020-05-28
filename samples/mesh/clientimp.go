package main

import (
	"context"
	"enceladus/common/xlog"
	mproto "enceladus/project/mesh_demo/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

func callEntrySvc(conn *grpc.ClientConn, svc, port string, cmd, input int) error {
	var respStr string
	var err error
	client := mproto.NewEntryClient(conn)
	switch cmd {
	case int(mproto.C2S_Cmds_Entry_Insert):
		in := &mproto.InsertReq{Val: int32(input)}
		resp, erra := client.Insert(context.Background(), in)
		respStr = resp.String()
		if len(respStr) < 1 {
			respStr = fmt.Sprintf("%+v", *resp)
		}
		err = erra
	case int(mproto.C2S_Cmds_Entry_Sum):
		in := &mproto.SumReq{}
		resp, erra := client.Sum(context.Background(), in)
		respStr = resp.String()
		if len(respStr) < 1 {
			respStr = fmt.Sprintf("%+v", *resp)
		}
		err = erra
	case int(mproto.C2S_Cmds_Entry_Max):
		in := &mproto.MaxReq{}
		resp, erra := client.Max(context.Background(), in)
		respStr = resp.String()
		if len(respStr) < 1 {
			respStr = fmt.Sprintf("%+v", *resp)
		}
		err = erra
	default:
		err = fmt.Errorf("invalid cmd: %v", cmd)
	}

	if err != nil {
		xlog.Error("callEntrySvc svc=%v port=%v cmd=%v input=%v err=%v", svc, port, cmd, input, err)
	} else {
		xlog.Info("callEntrySvc svc=%v port=%v cmd=%v input=%v resp=%v", svc, port, cmd, input, respStr)
	}
	return nil
}

func callMiddleSvc(conn *grpc.ClientConn, svc, port string, cmd, input int) error {
	var respStr string
	var err error
	client := mproto.NewMiddleClient(conn)
	switch cmd {
	case int(mproto.C2S_Cmds_Middle_InsertDecorator):
		in := &mproto.InsertReq{Val: int32(input)}
		resp, erra := client.InsertDecorator(context.Background(), in)
		respStr = resp.String()
		if len(respStr) < 1 {
			respStr = fmt.Sprintf("%+v", *resp)
		}
		err = erra
	default:
		err = fmt.Errorf("invalid cmd: %v", cmd)
	}
	if err != nil {
		xlog.Error("callEntrySvc svc=%v port=%v cmd=%v input=%v err=%v", svc, port, cmd, input, err)
	} else {
		xlog.Info("callEntrySvc svc=%v port=%v cmd=%v input=%v resp=%v", svc, port, cmd, input, respStr)
	}
	return nil
}

func callModelSvc(conn *grpc.ClientConn, svc, port string, cmd, input int) error {
	var resp proto.Message
	var err error
	client := mproto.NewModelClient(conn)
	switch cmd {
	case int(mproto.C2S_Cmds_Model_Insert):
		in := &mproto.InsertReq{Val: int32(input)}
		resp, err = client.Insert(context.Background(), in)
	case int(mproto.C2S_Cmds_Model_Len):
		in := &mproto.LenReq{}
		resp, err = client.Len(context.Background(), in)
	case int(mproto.C2S_Cmds_Model_Sum):
		in := &mproto.SumReq{}
		resp, err = client.Sum(context.Background(), in)
	case int(mproto.C2S_Cmds_Model_Max):
		in := &mproto.MaxReq{}
		resp, err = client.Max(context.Background(), in)
	default:
		err = fmt.Errorf("invalid cmd: %v", cmd)
	}
	respStr := ""
	if resp != nil {
		/*
			note:
				tmp := &mproto.SumResp{Sum: 0}
				str1 := proto.MarshalTextString(tmp)
				str2 := tmp.String()
				str1 和 str2 都是空字符串，除非tmp.Sum != 0
		*/
		respStr = resp.String()
	}
	if err != nil {
		xlog.Error("callEntrySvc svc=%v port=%v cmd=%v input=%v err=%v", svc, port, cmd, input, err)
	} else {
		xlog.Info("callEntrySvc svc=%v port=%v cmd=%v input=%v resp=%v", svc, port, cmd, input, respStr)
	}
	return nil
}

func clientFuncInsecure(svc, port string, cmd, input int) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		xlog.Error("c----- clientfunc.Dial: err %v", err)
		return
	}
	defer conn.Close()

	switch svc {
	case "entry":
		err = callEntrySvc(conn, svc, port, cmd, input)
	case "middle":
		err = callMiddleSvc(conn, svc, port, cmd, input)
	case "model":
		err = callModelSvc(conn, svc, port, cmd, input)
	default:
		err = fmt.Errorf("invalid svc: %v", svc)
	}
	if err != nil {
		xlog.Error("clientFuncInsecure svc=%v port=%v cmd=%v input=%v err=%v", err)
	}

	return
}
