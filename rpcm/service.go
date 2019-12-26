package rpcm

import (
	"github.com/joyous-x/saturn/gins"
	"github.com/joyous-x/saturn/rpcm/base"
)

func NewService(conf *base.ServiceConfig) (base.IMServer, error) {
	switch conf.Protocal {
	case "http":
		return newGinServer(conf)
	case "websocket":
	default:
	}
	return nil, nil
}

func newGinServer(conf *base.ServiceConfig) (base.IMServer, error) {
	server := gins.NewGinServer()
	err := server.Init(conf, gins.DefaultMiddlewares...)
	if err != nil {
		return nil, err
	}
	return server, err
}
