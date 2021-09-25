package mimic

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"net"
)

type dialer func(ctx context.Context, target string) (conn net.Conn, err error)

func dialAdapter(pconn *net.UDPConn) (dialer dialer) {

	dialer = func(ctx context.Context, target string) (conn net.Conn, err error) {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}

		//log.Printf("quic.DialAddrContext target=%v", target)
		//log.Printf("quic.DialAddrContext pconn=%v", pconn.LocalAddr())

		rAddr, err := net.ResolveUDPAddr("udp", target)
		if err != nil {
			return
		}

		//log.Printf("quic.DialAddrContext rAddr=%v", rAddr)

		sess, err := quic.DialContext(ctx, pconn, rAddr, "", tlsConf, config)
		if err != nil {
			return
		}

		//log.Printf("quic.DialAddrContext OpenStreamSync target=%v", target)
		stream, err := sess.OpenStreamSync(ctx)
		if err != nil {
			return
		}

		//log.Printf("connected target=%v", target)

		conn = &QUICConn{Sess: sess, Stream: stream}

		return
	}

	return
}
