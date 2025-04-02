package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIpAddress(request *http.Request) string {
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		forwardedFor := request.Header.Get("X-Forwarded-For")

		if forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")
			ip = parts[0]
		}
	}

	return ip
}
