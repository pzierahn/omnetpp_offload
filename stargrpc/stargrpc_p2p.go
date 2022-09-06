package stargrpc

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/mimic"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

// DialP2P creates a gRPC client connection over peer-to-peer.
func DialP2P(ctx context.Context, addr stargate.DialAddr) (cc *grpc.ClientConn, err error) {

	log.Printf("DialP2P: %v", addr)

	ctx, cln := context.WithTimeout(ctx, time.Second*5)
	defer cln()

	gate, raddr, err := stargate.DialP2PUDP(ctx, addr)
	if err != nil {
		return
	}

	adapter := mimic.NewDialAdapter(gate)

	return grpc.DialContext(
		ctx,
		raddr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(adapter),
	)
}

// ServeP2P establishes a peer-to-peer connection over stargate to serve the server.
func ServeP2P(addr stargate.DialAddr, server *grpc.Server) {
	ctx := context.Background()

	for {

		//
		// Establish a connection.
		//

		p2p, _, err := stargate.DialP2PUDP(ctx, addr)
		if err != nil {
			log.Println(err)
			continue
		}

		listener, err := mimic.NewQUICListener(p2p)
		if err != nil {
			log.Println(err)
			_ = p2p.Close()
			continue
		}

		log.Printf("ServeP2P: new connection addr=%v", listener.Addr())

		//
		// Detach serving process.
		//

		go func() {
			defer func() {
				_ = listener.Close()
				_ = p2p.Close()
			}()

			if err := server.Serve(listener); err != nil {
				log.Println(err)
			}
		}()
	}
}
