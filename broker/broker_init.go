package broker

import (
	"com.github.patrickz98.omnet/defines"
	pb "com.github.patrickz98.omnet/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Broker ", log.LstdFlags|log.Lshortfile)
}

type broker struct {
	pb.UnimplementedBrokerServer
	workers map[string]*pb.ClientInfo
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
		workers: make(map[string]*pb.ClientInfo),
		queue: queue{
			mu:      sync.Mutex{},
			jobs:    make(WorkHeap, 0),
			workers: make(map[string]chan *pb.Work),
		},
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)

	err = server.Serve(lis)

	return
}
