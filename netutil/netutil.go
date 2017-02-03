package netutil

import (
	"net"
	"time"
)

type TCPKeepAliveListener struct {
	*net.TCPListener
	duration time.Duration
}

func NewTCPKeepAliveListener(ln *net.TCPListener, d time.Duration) *TCPKeepAliveListener {
	return &TCPKeepAliveListener{
		TCPListener: ln,
		duration:    d,
	}
}

func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	if ln.duration == 0 {
		ln.duration = 3 * time.Minute
	}
	tc.SetKeepAlivePeriod(ln.duration)
	return tc, nil
}
