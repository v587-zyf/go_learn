package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetLocalIp() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func GetIpAddress(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(forwardedFor, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return ip
			}
		}
	}
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}
	index := strings.LastIndex(r.RemoteAddr, ":")
	if index < 0 {
		return r.RemoteAddr
	}
	return r.RemoteAddr[:index]
}
