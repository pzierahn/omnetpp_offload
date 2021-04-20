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
	workers map[string]*pb.ClientInfo
	work    chan *pb.Simulation
}

func Start() (err error) {

	var lis net.Listener
	lis, err = net.Listen("tcp", defines.Port)
	if err != nil {
		return
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &broker{
		workers: make(map[string]*pb.ClientInfo),
		work:    make(chan *pb.Simulation),
	})
	err = server.Serve(lis)

	return
}
