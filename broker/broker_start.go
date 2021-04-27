package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"net"
)

type broker struct {
	pb.UnimplementedBrokerServer
	workers workerList
	queue   queue
}

func Start(conf Config) (err error) {

	logger.Println("start server on", conf.Port)

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		return
	}

	defer func() { _ = lis.Close() }()

	brk := broker{
		workers: initWorkerList(),
		queue:   initQueue(),
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Server{})

	err = server.Serve(lis)

	return
}
