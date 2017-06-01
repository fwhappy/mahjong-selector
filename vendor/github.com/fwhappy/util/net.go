package util

import (
	"net"
	"strings"
)

// GetIP 读取客户端IP
func GetIP(conn *net.TCPConn) string {
	remote := conn.RemoteAddr().String()
	addrs := strings.Split(remote, ":")
	return addrs[0]
}
