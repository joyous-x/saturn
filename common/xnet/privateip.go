package xnet

import (
	"fmt"
	"net"
)

var privateCIDRs []*net.IPNet = privateCIDRs2IPNet()

func privateCIDRs2IPNet() []*net.IPNet {
	var cidrs = []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918, 10.0.0.0/8：10.0.0.0 ~ 10.255.255.255
		"172.16.0.0/12",  // RFC1918, 172.16.0.0/12：172.16.0.0 ~ 172.31.255.255
		"192.168.0.0/16", // RFC1918, 192.168.0.0/16：192.168.0.0 ~ 192.168.255.255
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	}

	ipnets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Sprintf("invalid CIDR address: in=%q err=%v", cidr, err))
		}
		ipnets = append(ipnets, ipnet)
	}
	return ipnets
}

// IsPrivateIP 比较高效的 private ip 判断
func IsPrivateIP(ip string) bool {
	return isPrivateIP(net.ParseIP(ip))
}

// IsPrivateIPEx 判断是否是 private ip
func IsPrivateIPEx(ip string) bool {
	return IsPrivateIPEx(net.ParseIP(ip))
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if nil == ip4 {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

func isPrivateIPEx(ip net.IP) bool {
	for _, ipnet := range privateCIDRs {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}
