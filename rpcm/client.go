package rpcm

import (
	"github.com/joyous-x/saturn/rpcm/base"
)

func GetClient(conf *base.ClientConfig) (base.IMClient, error) {
	switch conf.Protocal {
	case "http":
		return newGinClient(conf)
	case "websocket":
	default:
	}
	return nil, nil
}

type GinClient struct {
}

func (this *GinClient) Call(method, relativePath string, req, resp interface{}) error {
	return nil
}

func newGinClient(conf *base.ClientConfig) (base.IMClient, error) {
	return &GinClient{}, nil
}
