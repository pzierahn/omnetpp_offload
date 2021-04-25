package broker

import (
	"context"
	"encoding/json"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) Status(ctx context.Context, req *pb.StatusRequest) (reply *pb.StatusReply, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "  ")
	logger.Printf("Status: %server", jsonBytes)

	reply = &pb.StatusReply{}

	return
}
