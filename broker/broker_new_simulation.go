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
		for _, conf := range req.GetConfigs() {
			for _, run := range conf.RunNumbers {
				work := pb.Work{
					SimulationId: req.SimulationId,
					Source:       req.Source,
					Config:       conf.Name,
					RunNumber:    run,
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
