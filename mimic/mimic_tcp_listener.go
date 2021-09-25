package mimic

import (
	"fmt"
	"net"
)

type ListenerTCP struct {
	conn  *net.TCPConn
	block chan bool
}

// TCPConnToListener wraps an already established TCP connection into a TCP listener.
func TCPConnToListener(conn *net.TCPConn) (lis net.Listener) {
	ch := make(chan bool, 1)
	ch <- true

	lis = &ListenerTCP{
		conn:  conn,
		block: ch,
	}

	return
}

// Accept waits for and returns the next connection to the listener.
func (l *ListenerTCP) Accept() (conn net.Conn, err error) {
	if _, ok := <-l.block; ok {
		conn = l.conn
	} else {
		err = fmt.Errorf("connection cloesed")
	}

	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *ListenerTCP) Close() error {
	defer close(l.block)
	return l.conn.Close()
}

// Addr returns the listener's network address.
func (l *ListenerTCP) Addr() net.Addr {
	return l.conn.LocalAddr()
}
