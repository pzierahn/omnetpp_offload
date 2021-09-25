package mimic

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"net"
	"time"
)

var config = &quic.Config{
	KeepAlive:      true,
	MaxIdleTimeout: time.Millisecond * 2000,
}

type DialAdapter func(ctx context.Context, target string) (conn net.Conn, err error)

func NewDialAdapter(udpConn *net.UDPConn) (adapter DialAdapter) {

	adapter = func(ctx context.Context, target string) (conn net.Conn, err error) {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}

		rAddr, err := net.ResolveUDPAddr("udp", target)
		if err != nil {
			return
		}

		sess, err := quic.DialContext(ctx, udpConn, rAddr, "", tlsConf, config)
		if err != nil {
			return
		}

		stream, err := sess.OpenStreamSync(ctx)
		if err != nil {
			return
		}

		return &QUICConn{
			Sess:   sess,
			Stream: stream,
		}, nil
	}

	return
}

// NewQUICListener creates a QUIC connection from a UDP connection.
// It returns a QUICListener which implements the net.Listener interface.
func NewQUICListener(conn *net.UDPConn) (p2p net.Listener, err error) {

	tlsConf, _ := generateTLSConfig()

	qLis, err := quic.Listen(conn, tlsConf, config)
	if err != nil {
		return
	}

	return &QUICListener{ql: qLis}, err
}
