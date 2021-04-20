package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"encoding/json"
)

func (server *broker) NewSimulation(_ context.Context, req *pb.Simulation) (reply *pb.SimulationReply, err error) {

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(req, "", "    ")
	if err != nil {
		return
	}

	logger.Println("new simulation", string(jsonBytes))

	go func() {
		server.work <- req
	}()

	reply = &pb.SimulationReply{
		SimulationId: req.SimulationId,
	}

	return
}
