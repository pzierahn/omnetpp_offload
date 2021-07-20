package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/equic"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
)

func (prov *provider) listenRelay() {

	for {
		ctx := context.Background()
		conn, err := stargate.RelayDialTCP(ctx, prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			log.Printf("listenRelay: listening LocalAddr=%v RemoteAddr=%v",
				conn.LocalAddr(), conn.RemoteAddr())

			server := grpc.NewServer()
			pb.RegisterProviderServer(server, prov)
			pb.RegisterStorageServer(server, prov.store)
			err := server.Serve(equic.ListenTCP(conn))
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
