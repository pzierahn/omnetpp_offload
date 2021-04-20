package broker

import (
	"com.github.patrickz98.omnet/defines"
	pb "com.github.patrickz98.omnet/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Broker ", log.LstdFlags|log.Lshortfile)
}

type broker struct {
	pb.UnimplementedBrokerServer
	workers workerList
	queue   queue
}

func Start() (err error) {

	logger.Println("start server")

	var lis net.Listener
	lis, err = net.Listen("tcp", defines.Port)
	if err != nil {
		return
	}

	brk := broker{
		workers: initWorkerList(),
		queue:   initQueue(),
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)

	err = server.Serve(lis)

	return
}
