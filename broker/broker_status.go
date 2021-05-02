package broker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) Status(_ context.Context, req *pb.StatusRequest) (reply *pb.StatusReply, err error) {

	logger.Printf("status %v\n", req.SimulationIds)
	reply, err = server.db.Status(req)

	return
}
