package broker

import (
	"context"
	"encoding/json"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) Push(ctx context.Context, req *pb.WorkResult) (reply *pb.WorkAffirmation, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "  ")
	logger.Println("results", string(jsonBytes))

	reply = &pb.WorkAffirmation{}

	return
}
