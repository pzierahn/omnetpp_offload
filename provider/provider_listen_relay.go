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
	relay := pb.NewStargateClient(bconn)

	for {
		port, err := relay.Relay(context.Background(), &pb.RelayRequest{
			DialAddr: prov.providerId,
		})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Connect over relay server (port: %v)", port.Port)

		// TODO: Replace gconfig.Config.Broker.Address
		raddr := &net.TCPAddr{
			IP:   net.ParseIP(gconfig.Config.Broker.Address),
			Port: int(port.Port),
		}

		conn, err := net.DialTCP("tcp", &net.TCPAddr{}, raddr)
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			log.Printf("########## start listening LocalAddr=%v RemoteAddr=%v",
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
