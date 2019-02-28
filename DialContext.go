package net

import (
	"context"
	traditionalnet "net"
	"os"
)

// Dial makes TCP connection.
func Dial(network, remote string) (conn traditionalnet.Conn, err error) {
	httpproxy := os.Getenv(`HTTPS_PROXY`)
	if httpproxy != `` {
		conn, err = proxyDial(network, remote, httpproxy)
	}
	if conn == nil {
		conn, err = traditionalnet.Dial(network, remote)
	}
	return conn, err
}

// DialContext makes TCP connection.
func DialContext(ctx context.Context, network, remote string) (traditionalnet.Conn, error) {
	httpproxy := os.Getenv(`HTTPS_PROXY`)
	if httpproxy != `` {
		return proxyDialContext(ctx, network, remote, httpproxy)
	} else {
		var dialer traditionalnet.Dialer
		return dialer.DialContext(ctx, network, remote)
	}

}
