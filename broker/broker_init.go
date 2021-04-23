package broker

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
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
