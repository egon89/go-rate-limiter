package utils

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIpAddress(request *http.Request) (string, error) {
	ip := GetIpAddressFromXForwardedFor(request)
	if ip != "" {
		return ip, nil
	}

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get ip address: %w", err)
	}

	return ip, nil
}

func GetIpAddressFromXForwardedFor(request *http.Request) string {
	forwardedFor := request.Header.Get("X-Forwarded-For")
	if forwardedFor == "" {
		return ""
	}

	parts := strings.Split(forwardedFor, ",")

	return strings.TrimSpace(parts[0])
}
