package provider

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (prov *provider) listenP2P() {
	for {
		log.Println("wait for peer to peer connect")

		p2p, err := stargate.DialQUICgRPCListener(context.Background(), prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		go func(p2p net.Listener) {

			defer func() { _ = p2p.Close() }()

			server := grpc.NewServer()
			pb.RegisterProviderServer(server, prov)
			pb.RegisterStorageServer(server, prov.store)
			err := server.Serve(p2p)
			if err != nil {
				log.Println(err)
			}
		}(p2p)
	}
}
