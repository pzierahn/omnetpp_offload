package broker

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) ExecuteSimulation(_ context.Context, req *pb.Simulation) (reply *pb.SimulationReply, err error) {

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(req, "", "  ")
	if err != nil {
		return
	}

	logger.Println("execute simulation", string(jsonBytes))

	if req.SimulationId == "" {
		err = fmt.Errorf("SimulationId missing")
		return
	}

	go func() {
		server.db.NewSimulation(req)
		server.db.DistributeWork()
	}()

	reply = &pb.SimulationReply{
		SimulationId: req.SimulationId,
	}

	return
}
