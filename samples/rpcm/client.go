package main

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/rpcm"
	"github.com/joyous-x/saturn/rpcm/base"
)

func main() {
	xlog.Debug("rpcm.HTTP sample ===> start ")
	conf := &base.ClientConfig{
		Protocal: "http",
		Encoding: "json",
		Scheme:   "http",
		Name:     "",
		Address:  "127.0.0.1:8001",
	}
	iclient, err := rpcm.GetClient(conf)
	if err != nil {
		xlog.Error("get client err:%v config:%+v", err, *conf)
		return
	}
	req := map[string]string{
		"hi": "hello",
	}
	resp := make(map[string]string, 0)
	err = iclient.Call("POST", "/v1", req, &resp)
	if err != nil {
		xlog.Error("client.Call(%v %v) err:%v", "POST", "/v1", err)
		return
	}
	xlog.Debug("client.Call(%v %v) resp:%+v", "POST", "/v1", resp)
	xlog.Debug("rpcm.HTTP sample ===> end ")
}
