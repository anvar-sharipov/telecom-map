package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	// 1. Если есть proxy
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// XFF может быть "client, proxy1, proxy2"
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// 2. Без proxy
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
