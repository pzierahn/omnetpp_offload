package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
)

const (
	port = ":50051"
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

func (server *broker) WorkFinished(ctx context.Context, req *pb.WorkResult) (reply *pb.WorkAffirmation, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "    ")
	logger.Printf("Status: %server", jsonBytes)

	reply = &pb.WorkAffirmation{}

	return
}

func (server *broker) PinPong(stream pb.Broker_PinPongServer) (err error) {

	inx := 0

	for {
		message, err := stream.Recv()
		if err != nil {
			logger.Println(err)
			break
		}

		logger.Println("Receive", message.Message, message.Time)

		if inx == 0 {
			logger.Println("Send reply", message.Message, message.Time)
			err = stream.Send(&pb.Pong{
				Message: "Pong",
				Time:    timestamppb.Now(),
			})
			if err != nil {
				logger.Println(err)
			}
		}

		inx = (inx + 1) % 10
	}

	return
}

func Start() (err error) {

	var lis net.Listener
	lis, err = net.Listen("tcp", port)
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
