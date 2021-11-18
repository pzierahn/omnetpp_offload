package broker

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/stargate"
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
	pb.RegisterEvaluationServer(server, &eval.Server{})
	err = server.Serve(lis)

	return
}
