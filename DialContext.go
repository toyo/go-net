package net

import (
	"context"
	traditionalnet "net"
	"os"
)

// Dial makes TCP connection.
func Dial(network, remote string) (conn *traditionalnet.TCPConn, err error) {
	httpproxy := os.Getenv(`HTTPS_PROXY`)
	if httpproxy != `` {
		return proxyDial(network, remote, httpproxy)
	}
	tcpaddr, _ := traditionalnet.ResolveTCPAddr(network, remote)
	return traditionalnet.DialTCP(network, nil, tcpaddr)
}

// DialContext makes TCP connection.
func DialContext(ctx context.Context, network, remote string) (conn *traditionalnet.TCPConn, err error) {
	connch := make(chan *traditionalnet.TCPConn)
	errch := make(chan error)
	go func() {
		conn, err := Dial(network, remote)
		connch <- conn
		errch <- err
	}()

	select {
	case conn = <-connch:
		err = <-errch
		return
	case <-ctx.Done():
		remoteaddr, _ := traditionalnet.ResolveTCPAddr(network, remote)
		err = &traditionalnet.OpError{
			Op:     `connect`,
			Net:    network,
			Err:    ctx.Err(),
			Addr:   remoteaddr,
			Source: nil}
		return
	}
}
