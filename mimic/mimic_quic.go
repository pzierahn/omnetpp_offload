package mimic

import (
	"context"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

type QUICConn struct {
	Sess   quic.Connection
	Stream quic.Stream
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *QUICConn) Read(b []byte) (n int, err error) {
	return c.Stream.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *QUICConn) Write(b []byte) (n int, err error) {
	return c.Stream.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *QUICConn) Close() error {
	return c.Stream.Close()
}

// LocalAddr returns the local network address.
func (c *QUICConn) LocalAddr() net.Addr {
	return c.Sess.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *QUICConn) RemoteAddr() net.Addr {
	return c.Sess.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *QUICConn) SetDeadline(t time.Time) error {
	return c.Stream.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *QUICConn) SetReadDeadline(t time.Time) error {
	return c.Stream.SetReadDeadline(t)

}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *QUICConn) SetWriteDeadline(t time.Time) error {
	return c.Stream.SetWriteDeadline(t)
}

type QUICListener struct {
	ql quic.Listener
}

// Accept waits for and returns the next connection to the listener.
func (l *QUICListener) Accept() (net.Conn, error) {
	sess, err := l.ql.Accept(context.Background())
	if err != nil {
		return nil, err
	}

	s, err := sess.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}

	return &QUICConn{sess, s}, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *QUICListener) Close() error {
	return l.ql.Close()
}

// Addr returns the listener's network address.
func (l *QUICListener) Addr() net.Addr {
	return l.ql.Addr()
}
