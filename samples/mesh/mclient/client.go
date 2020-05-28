package mclient

import (
	"fmt"
)

//>  ip:port server_name svc_name

type ServerInfo struct {
	Host string
	Port int32
	Name string
}

var GrpcServers map[string]*ServerInfo

func RegistGrpcServer(name, host string) error {
	if _, ok := GrpcServers[name]; ok {
		return fmt.Errorf("server %v already exists", name)
	}
	return nil
}
