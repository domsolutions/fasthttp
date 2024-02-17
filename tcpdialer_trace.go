package fasthttp

import (
	"net"
	"sync/atomic"
)

type TCPDialerTrace struct {
	td           TCPDialer
	BytesRead    int64
	BytesWritten int64
}

func (d *TCPDialerTrace) Dial(addr string) (net.Conn, error) {
	conn, err := d.td.dial(addr, false, DefaultDialTimeout)
	if err != nil {
		return nil, err
	}

	return &ConnTrace{conn, &d.BytesWritten, &d.BytesRead}, nil
}

func NewTCPDialTrace() *TCPDialerTrace {
	return &TCPDialerTrace{td: TCPDialer{Concurrency: 1000}}
}

func DialTrace(d *TCPDialerTrace) DialFunc {
	return func(addr string) (net.Conn, error) {
		return d.Dial(addr)
	}
}

type ConnTrace struct {
	net.Conn
	BytesWritten *int64
	BytesRead    *int64
}

func (c *ConnTrace) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if n > 0 {
		atomic.AddInt64(c.BytesRead, int64(n))
	}
	return n, err
}

func (c *ConnTrace) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if n > 0 {
		atomic.AddInt64(c.BytesWritten, int64(n))
	}
	return n, err
}
