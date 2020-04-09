package xnet

import (
	"net"
	"strconv"
)

// GetLocalIdlePort ...
func GetLocalIdlePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	_, portString, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portString)
}
