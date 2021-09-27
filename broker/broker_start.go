package broker

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/eval"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start(config gconfig.Broker) (err error) {

	stargate.SetConfig(stargate.Config{
		Addr: config.Address,
		Port: config.StargatePort,
	})

	go func() {
		err := stargate.Server(context.Background(), true)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Printf("start broker on :%d", config.BrokerPort)

	brk := broker{
		providers:   make(map[string]*pb.ProviderInfo),
		utilization: make(map[string]*pb.Utilization),
		listener:    make(map[string]chan<- *pb.ProviderList),
	}

	go brk.startDebugWebAPI()

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.BrokerPort,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterEvalServer(server, &eval.Server{})
	err = server.Serve(lis)

	return
}
