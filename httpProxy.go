package net

import (
	"bufio"
	"context"
	traditionalnet "net"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

// httpProxy is a HTTP/HTTPS proxy class.
type httpProxy struct {
	host     string
	needAuth bool
	user     string
	pass     string
	dialer   proxy.Dialer
}

func init() {
	proxy.RegisterDialerType("http", newHTTPProxy)
	proxy.RegisterDialerType("https", newHTTPProxy)
}

func newHTTPProxy(uri *url.URL, dialer proxy.Dialer) (proxy.Dialer, error) {
	s := new(httpProxy)
	s.host = uri.Host
	s.dialer = dialer
	if uri.User != nil {
		s.needAuth = true
		s.user = uri.User.Username()
		s.pass, _ = uri.User.Password()
	}

	return s, nil
}

func (s *httpProxy) Dial(network, addr string) (traditionalnet.Conn, error) {
	if network != `tcp` {
		return nil, errors.New(`Only TCP supported`)
	}
	proxyconn, err := s.dialer.Dial(network, s.host)
	if err != nil {
		return nil, err
	}

	remoteurl, err := url.Parse("http://" + addr)
	if err != nil {
		proxyconn.Close()
		return nil, err
	}
	remoteurl.Scheme = ""

	proxyreq, err := http.NewRequest("CONNECT", remoteurl.String(), nil)
	if err != nil {
		proxyconn.Close()
		return nil, err
	}
	proxyreq.Close = false
	if s.needAuth {
		proxyreq.SetBasicAuth(s.user, s.pass)
	}
	proxyreq.Header.Set("User-Agent", "Golang proxy agent")

	if err = proxyreq.Write(proxyconn); err != nil {
		proxyconn.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(proxyconn), proxyreq)
	if err != nil {
		resp.Body.Close()
		proxyconn.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		proxyconn.Close()
		err = errors.New("Proxy server error: " + resp.Status)
		return nil, err
	}

	return proxyconn, nil
}

func proxyDial(network, remote, httpproxy string) (conn traditionalnet.Conn, err error) {
	var u *url.URL
	u, err = url.Parse(httpproxy)
	if err == nil {
		var dialer traditionalnet.Dialer
		var d proxy.Dialer
		d, err = proxy.FromURL(u, &dialer)
		if err == nil {
			conn, err = d.Dial(network, remote)
		}
	}
	return
}

func proxyDialContext(ctx context.Context, network, remote, httpproxy string) (conn traditionalnet.Conn, err error) {
	connch := make(chan traditionalnet.Conn)
	errch := make(chan error)
	go func() {
		conn, err := proxyDial(network, remote, httpproxy)
		connch <- conn
		errch <- err
	}()

	select {
	case conn = <-connch:
		err = <-errch
		return
	case <-ctx.Done():
		remoteaddr, _ := traditionalnet.ResolveTCPAddr(network, remote)
		proxyaddr, _ := traditionalnet.ResolveTCPAddr(network, httpproxy)
		err = &traditionalnet.OpError{
			Op:     `connect`,
			Net:    network,
			Err:    ctx.Err(),
			Addr:   remoteaddr,
			Source: proxyaddr}
		return
	}
}
