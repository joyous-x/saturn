package mserver

import (
	"context"
	wproto "enceladus/wgrpc/protoc"
)

// NewDispatchResp ...
func NewDispatchResp(ctx context.Context, comm *wproto.RouteHeader) (*wproto.DispatchResp, error) {
	resp := &wproto.DispatchResp{
		Header: &wproto.RouteHeader{
			ReqId: comm.ReqId,
			Uid:   comm.Uid,
			Appid: comm.Appid,
			Cmd:   comm.Cmd,
		},
		Cmd:     comm.Cmd + 1,
		RetCode: -1,
	}
	return resp, nil
}
