package broker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) Status(_ context.Context, req *pb.ResultsRequest) (reply *pb.StatusReply, err error) {

	logger.Println("status", req.SimulationId)
	reply, err = server.db.Status(req)

	return
}
