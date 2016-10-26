package utils

import (
	"net"
)

func createTCPListener(host, port string) (*net.TCPListener, error) {
	addr := net.JoinHostPort(host, port)
	if tcpAddr, err := net.ResolveTCPAddr("tcp4", addr); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp4", tcpAddr)
	}
}
