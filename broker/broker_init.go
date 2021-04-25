package broker

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
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

	defer func() { _ = lis.Close() }()

	brk := broker{
		workers: initWorkerList(),
		queue:   initQueue(),
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Storage{})

	err = server.Serve(lis)

	return
}
