package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/equic"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (prov *provider) listenRelay(bconn *grpc.ClientConn) {

	//
	// TODO: remove log.Fatalln()
	//

	relay := pb.NewStargateClient(bconn)

	for {
		port, err := relay.Relay(context.Background(), &pb.RelayRequest{
			DialAddr: prov.providerId,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("listenRelay: port=%v", port.Port)

		// TODO: Replace this with stargate.Dial...
		raddr, err := net.ResolveTCPAddr("tcp", gconfig.BrokerDialAddr())
		if err != nil {
			log.Fatalln(err)
		}

		raddr.Port = int(port.Port)

		log.Printf("listenRelay: relay=%v", raddr.String())

		conn, err := net.DialTCP("tcp", &net.TCPAddr{}, raddr)
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
