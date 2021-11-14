package stargrpc

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/mimic"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// DialRelay creates a gRPC client connection over the stargate relay server to the dial address.
func DialRelay(ctx context.Context, addr stargate.DialAddr) (cc *grpc.ClientConn, err error) {

	log.Printf("DialRelay: addr=%v", addr)

	ctx, cln := context.WithTimeout(ctx, time.Second*5)
	defer cln()

	conn, err := stargate.DialRelayTCP(ctx, addr)
	if err != nil {
		return
	}

	log.Printf("DialRelay: addr=%v remote=%v", addr, conn.RemoteAddr().String())

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

// ServeRelay establishes a relay connection over stargate to serve the server.
func ServeRelay(addr stargate.DialAddr, server *grpc.Server) {
	ctx := context.Background()

	for {

		//
		// Establish a connection.
		//

		conn, err := stargate.DialRelayTCP(ctx, addr)
		if err != nil {
			log.Println(err)
			continue
		}

		listener := mimic.TCPConnToListener(conn)

		log.Printf("ServeRelay: new connection addr=%v", conn.RemoteAddr())

		//
		// Detach serving process.
		//

		go func() {
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
