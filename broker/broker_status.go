package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"encoding/json"
)

func (server *broker) Status(ctx context.Context, req *pb.StatusRequest) (reply *pb.StatusReply, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "    ")
	logger.Printf("Status: %server", jsonBytes)

	reply = &pb.StatusReply{}

	return
}
