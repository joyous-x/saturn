package rpcm

import (
	"encoding/json"
	"fmt"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/rpcm/base"
	"net/url"
	"strings"
)

func GetClient(conf *base.ClientConfig) (base.IMClient, error) {
	switch strings.ToLower(conf.Protocal) {
	case "http":
		return newGinClient(conf)
	case "websocket":
	default:
	}
	return nil, nil
}

type GinClient struct {
	conf *base.ClientConfig
}

func (this *GinClient) Call(method, relativePath string, req, resp interface{}) error {
	// TODO, 暂时只支持 json 格式
	if this.conf.Encoding != "json" {
		return fmt.Errorf("invalid encoding: %v", this.conf.Encoding)
	}
	path := &url.URL{
		Scheme: this.conf.Protocal,
		Host:   this.conf.Address,
		Path:   relativePath,
	}
	datas, err := utils.HttpPostJson(path.String(), req)
	if err == nil {
		err = json.Unmarshal(datas, resp)
	}
	if err != nil {
		xlog.Error("GinClient.Call(%v) error:%v", path.String(), err)
	}
	return err
}

func newGinClient(conf *base.ClientConfig) (base.IMClient, error) {
	return &GinClient{
		conf: conf,
	}, nil
}
