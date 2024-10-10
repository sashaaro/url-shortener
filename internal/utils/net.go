package utils

import (
	"errors"
	"net"
	"os"
)

// GetMyHostIP - получить ip своего хоста
func GetMyHostIP() (net.IP, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	if len(addr) == 0 {
		return nil, errors.New("not found")
	}
	return addr[0], nil
}
