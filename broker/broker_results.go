package broker

import (
	"context"
	"encoding/json"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) PutResults(_ context.Context, req *pb.TaskResult) (reply *pb.WorkAffirmation, err error) {
	jsonBytes, _ := json.MarshalIndent(req, "", "  ")
	logger.Println("commit results", string(jsonBytes))

	server.db.ReceiveResult(req)

	reply = &pb.WorkAffirmation{}

	return
}
