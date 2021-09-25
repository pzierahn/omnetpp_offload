package egrpc

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/mimic"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"time"
)

func DialP2P(ctx context.Context, addr string) (cc *grpc.ClientConn, err error) {

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
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(adapter),
	)
}

func ServeP2P(addr string, server *grpc.Server) {
	ctx := context.Background()

	for {
		p2p, _, err := stargate.DialP2PUDP(ctx, addr)
		if err != nil {
			log.Println(err)
			continue
		}

		listener, err := mimic.NewQUICListener(p2p)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("ServeP2P: new connection addr=%v", listener.Addr())

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
