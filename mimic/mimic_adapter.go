package mimic

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"google.golang.org/grpc"
	"net"
)

func DialGRPC(ctx context.Context, remote string, udpConn *net.UDPConn) (conn *grpc.ClientConn, err error) {

	dialer := func(ctx context.Context, target string) (conn net.Conn, err error) {
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

		conn = &QUICConn{Sess: sess, Stream: stream}

		return
	}

	return grpc.DialContext(
		ctx,
		remote,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
	)
}

func NewQUICListener(conn *net.UDPConn) (p2p net.Listener, err error) {

	tlsConf, _ := generateTLSConfig()

	qLis, err := quic.Listen(conn, tlsConf, config)
	if err != nil {
		return
	}

	return &QUICListener{ql: qLis}, err
}
