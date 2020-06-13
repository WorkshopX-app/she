package httpx

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

var _private []*net.IPNet

func init() {
	spec := []string{
		//ipv4
		"127.0.0.1/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		//ipv6
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	_private = make([]*net.IPNet, len(spec))
	for pos, block := range spec {
		_, ipNet, _ := net.ParseCIDR(block)
		_private[pos] = ipNet
	}
}

//IsPrivateIP reports whether addr is private ip
//https://en.wikipedia.org/wiki/Private_network
func IsPrivateIP(addr string) (bool, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false, fmt.Errorf("addr:%s invalid", addr)
	}

	for pos := range _private {
		if _private[pos].Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

//RealIPFrom figure out the real ip from http request
// nginx X-Real-IP supported
// cloudflare true-client-ip supported
// cf > xForwardedFor > xRealIP > raw RemoteAddr
func RealIPFrom(r *http.Request) string {
	realIP := r.RemoteAddr
	//https://support.cloudflare.com/hc/en-us/articles/206776727-What-is-True-Client-IP-
	cfConnectingIP := strings.TrimSpace(r.Header.Get("CF-Connecting-IP"))
	//http://nginx.org/en/docs/http/ngx_http_realip_module.html
	xRealIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	xForwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))

	if cfConnectingIP != "" {

		return cfConnectingIP
	}

	if xForwardedFor != "" {
		for _, addr := range strings.Split(xForwardedFor, ",") {
			addr = strings.TrimSpace(addr)
			yes, err := IsPrivateIP(addr)
			if err == nil && !yes {
				return addr
			}
		}
	}

	if xRealIP != "" {
		return xRealIP
	}
	//both nothing, pop remote addr
	if strings.ContainsRune(realIP, ':') {
		realIP, _, _ = net.SplitHostPort(realIP)
	}
	return realIP
}
