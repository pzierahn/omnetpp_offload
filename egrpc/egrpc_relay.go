package egrpc

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/mimic"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func DialRelay(ctx context.Context, addr string) (cc *grpc.ClientConn, err error) {

	log.Printf("DialRelay: %v", addr)

	ctx, cln := context.WithTimeout(ctx, time.Second*5)
	defer cln()

	conn, err := stargate.DialRelayTCP(ctx, addr)
	if err != nil {
		return
	}

	log.Printf("DialRelay: dial %v", conn.RemoteAddr().String())

	return grpc.DialContext(
		ctx,
		conn.RemoteAddr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return conn, nil
		}),
	)
}

func ServeRelay(addr string, server *grpc.Server) {
	ctx := context.Background()

	for {
		conn, err := stargate.DialRelayTCP(ctx, addr)
		if err != nil {
			log.Println(err)
			continue
		}

		listener := mimic.TCPConnToListener(conn)

		go func() {
			log.Printf("serveRelay: listening LocalAddr=%v RemoteAddr=%v",
				conn.LocalAddr(), conn.RemoteAddr())

			defer func() {
				_ = listener.Close()
				_ = conn.Close()
			}()

			if err := server.Serve(listener); err != nil {
				log.Println(err)
			}
		}()
	}
}
