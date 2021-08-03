package broker

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start() (err error) {

	go func() {
		err := stargate.Server(context.Background(), true)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Printf("start broker on :%d", gconfig.BrokerPort())

	brk := broker{
		providers:   make(map[string]*pb.ProviderInfo),
		utilization: make(map[string]*pb.Utilization),
		listener:    make(map[string]chan<- *pb.Providers),
	}

	go brk.startWebService()

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: gconfig.BrokerPort(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Server{})
	err = server.Serve(lis)

	return
}
