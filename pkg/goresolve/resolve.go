// Package goresolve Functions that resolve IP addresses and search for host name.
// Caching mechanism, goroutine safe
package goresolve

import (
	"fmt"
	"net"
	"sync"
)

var (
	resolvedIp     = make(map[string][]string, 0)
	resolveIpMutex sync.Mutex
)

// Ip Resolve IP address - trying to find matching host name
func Ip(ip string) []string {
	resolveIpMutex.Lock()
	defer resolveIpMutex.Unlock()

	if resolvedIp[ip] != nil {
		return resolvedIp[ip]
	}

	ips, ok := net.LookupAddr(ip)
	if ok != nil {
		fmt.Println(ok)
		ips = []string{ip}
	}

	resolvedIp[ip] = ips
	return ips
}
