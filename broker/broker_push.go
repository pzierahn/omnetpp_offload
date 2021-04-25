package broker

import (
	"context"
	"encoding/json"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) CommitResults(ctx context.Context, req *pb.TaskResult) (reply *pb.WorkAffirmation, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "  ")
	logger.Println("commit results", string(jsonBytes))

	reply = &pb.WorkAffirmation{}

	return
}
