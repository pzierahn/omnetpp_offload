package equic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	"math/big"
	"net"
)

func GenerateTLSConfig() (tlsConf *tls.Config, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return
	}

	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}

	return
}

func GRPCDialerAuto() (udpconn *net.UDPConn, dialer Dialer) {

	udpconn, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	dialer = GRPCDialer(udpconn)

	return
}

type Dialer func(ctx context.Context, target string) (conn net.Conn, err error)

func GRPCDialer(pconn *net.UDPConn) (dialer Dialer) {

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

		sess, err := quic.DialContext(ctx, pconn, rAddr, "", tlsConf, &quic.Config{
			KeepAlive: true,
			//MaxIdleTimeout: time.Millisecond * 3200,
		})
		if err != nil {
			return
		}

		//log.Printf("quic.DialAddrContext OpenStreamSync target=%v", target)
		stream, err := sess.OpenStreamSync(ctx)
		if err != nil {
			return
		}

		//log.Printf("connected target=%v", target)

		conn = &Conn{Sess: sess, Stream: stream}

		return
	}

	return
}
