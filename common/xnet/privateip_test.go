package xnet 

import (
	"testing"
)

var testIPs = map[string]bool {
	"172.14.0.1": false,
	"1111.1401": false,
	"172.16.0.1": true,
	"172.18.0.1": true,
	"192.168.100.10": true,
}

func Test_IsPrivateIp(t *testing.T) {
	for ip,v := range testIPs {
		if v == IsPrivateIp(ip) {
			t.Logf("%v %v", ip, v)
		} else {
			t.Errorf("%v expect %v", ip, v)
		}
	}
}

func Test_IsPrivateIpEx(t *testing.T) {
	for ip,v := range testIPs {
		if v == IsPrivateIpEx(ip) {
			t.Logf("%v %v", ip, v)
		} else {
			t.Errorf("%v expect %v", ip, v)
		}
	}
}