package xnet 

import (
	"net"
	"net/http"
	"strings"
)

// X-Forwarded-For, X-Real-IP, remote_addr是http协议中用来表示客户端地址的请求头。
// X-Forwarded-For 和 X-Real-IP 只有请求存在代理\反向代理时才有值，而remote_addr一直存在。
//     X-Forwarded-For：记录代理服务器的地址。就像是vector中的push_back，每经过一个代理服务器就把该请求的来源地址加在记录的后面
//                      (由于是记录来源地址，所以该字段不会保存最后一个代理服务器的地址)。
//                      格式形如：1.1.1.1, 2.2.2.2
//     X-Real-IP：也是用来记录服务器的地址，但是和上面的不同，它不把记录添加到结尾，而是直接替换
//     remote_addr：上一个客户端连接的地址，不存在代理就表示客户端的地址，存在代理就表示最后一个代理服务器的地址
// note:
//     X-Forwarded-For和X-Real-IP 可以被客户端伪造，而remote_addr不能：
//         因为remote_addr字段不是通过请求头来决定的，而是服务端在建立tcp连接时获取的的客户端地址。
// eg:
//     请求：client -> client_proxy_A -> client_proxy_B -> server_Public 时，
//         X-Forwarded-For: client_IP, client_proxy_A (注意：经过了两次转发，但并没有记录 client_proxy_B 的地址。)
//         X-Real-IP: client_proxy_A (注意：会在请求经过 client_proxy_B 时被设置为 client_proxy_A 的地址)
//         remote_addr: client_proxy_B (or client_proxy_B_PublicIP 如果client_proxy_B为client内网代理时)


type IRealIP interface {
	IsPrivateIP(ip net.IP) bool
	RealIP(r *http.Request) string 
	ClientIP(r *http.Request) string 
	ClientPublicIP(r *http.Request) string 
}

type HttpRealIP struct {
}

// RealIP try to get the client's public address ip, if not, get the private one
func (m *HttpRealIP) RealIP(r *http.Request) string {
	realIp := m.ClientPublicIP(r)
	if realIp == ""{
		realIp = m.ClientIP(r)
	}
	return realIp
}

// IsPrivateIP check the ip address whether priviate or not
func (m *HttpRealIP) IsPrivateIP(ip string) bool {
	return m.isPrivateIP(net.ParseIP(ip))
}

// RemoteIP get remote ip from http.Request
func (m *HttpRealIP) RemoteIP(r *http.Request) string {
	remoteIP := ""
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		remoteIP = ip
	}
	return remoteIP
}

// ClientIP get client ip
func (m *HttpRealIP) ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// ClientPublicIP get client's public ip addr
func (m *HttpRealIP) ClientPublicIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	for _, tmpIP := range strings.Split(xForwardedFor, ",") {
		tmpIP = strings.TrimSpace(tmpIP)
		if tmpIP != "" && !m.IsPrivateIP(tmpIP) {
			return tmpIP
		}
	}
	xRealIp := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if xRealIp != "" && !m.IsPrivateIP(xRealIp) {
		return xRealIp
	}
	remoteIp, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && m.IsPrivateIP(remoteIp) {
		return remoteIp
	}
	return ""
}

// IsPrivateIP 
//     tcp/ip协议中，专门保留了三个IP地址区域作为私有地址，其地址范围如下：
//         10.0.0.0/8：10.0.0.0～10.255.255.255
//         172.16.0.0/12：172.16.0.0～172.31.255.255
//         192.168.0.0/16：192.168.0.0～192.168.255.255
func (m *HttpRealIP) isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4() 
	if nil == ip4 {
		return false
	}

	return ip4[0] == 10 ||                                 // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168)    // 192.168.0.0/16
}