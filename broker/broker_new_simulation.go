package broker

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) ExecuteSimulation(_ context.Context, req *pb.Simulation) (reply *pb.SimulationReply, err error) {

	// Todo: synchronize queue!

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(req, "", "    ")
	if err != nil {
		return
	}

	logger.Println("execute simulation", string(jsonBytes))

	if req.SimulationId == "" {
		err = fmt.Errorf("SimulationId missing")
		return
	}

	go func() {
		for _, run := range req.Run {
			for _, runNum := range run.RunNumbers {
				work := pb.Task{
					SimulationId: req.SimulationId,
					Simulation:   req.OppConfig,
					Source:       req.Source,
					Config:       run.Config,
					RunNumber:    runNum,
				}

				server.queue.jobs.Push(&work)
			}
		}

		server.distributeWork()
	}()

	reply = &pb.SimulationReply{
		SimulationId: req.SimulationId,
	}

	return
}
