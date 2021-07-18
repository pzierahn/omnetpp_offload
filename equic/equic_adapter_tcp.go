package equic

import (
	"log"
	"net"
)

type ListenerTCP struct {
	conn *net.TCPConn
	lock chan bool
}

func ListenTCP(conn *net.TCPConn) (lis net.Listener) {
	ch := make(chan bool, 1)
	ch <- true

	lis = &ListenerTCP{
		conn: conn,
		lock: ch,
	}

	return
}

// Accept waits for and returns the next connection to the listener.
func (l *ListenerTCP) Accept() (conn net.Conn, err error) {
	<-l.lock

	log.Printf("############# Accept")
	return l.conn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *ListenerTCP) Close() error {
	return l.conn.Close()
}

// Addr returns the listener's network address.
func (l *ListenerTCP) Addr() net.Addr {
	return l.conn.LocalAddr()
}
