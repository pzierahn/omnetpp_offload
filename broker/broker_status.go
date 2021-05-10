package broker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) SimulationStatus(_ context.Context, req *pb.StatusRequest) (reply *pb.StatusReply, err error) {

	logger.Printf("status %v\n", req.SimulationIds)
	reply, err = server.db.Status(req)

	return
}

func (server *broker) WorkerInfo(_ context.Context, _ *pb.WorkerInfoRequest) (reply *pb.WorkerInfoReply, err error) {

	logger.Printf("workerInfo")

	reply = &pb.WorkerInfoReply{}

	for _, info := range server.db.workerInfo {
		reply.Items = append(reply.Items, &pb.WorkerInfoReply_WorkerInfo{
			WorkerId:     info.WorkerId,
			Os:           info.Os,
			Arch:         info.Arch,
			NumCPUs:      uint32(info.NumCPUs),
			FreeResource: uint32(server.db.capacities[info.WorkerId]),
		})
	}

	return
}
