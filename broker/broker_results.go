package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"encoding/json"
)

func (server *broker) Results(ctx context.Context, req *pb.WorkResult) (reply *pb.WorkAffirmation, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "    ")
	logger.Println("results", string(jsonBytes))

	reply = &pb.WorkAffirmation{}

	return
}
