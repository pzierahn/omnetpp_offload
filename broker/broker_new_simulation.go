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
				confId := conf.Name + "." + run

				work := pb.Work{
					SimulationId: req.SimulationId,
					ConfigId:     confId,
					Source:       req.Source,
					Config:       conf.Name,
					RunNumber:    run,
				}

				logger.Println(work.SimulationId, work.ConfigId)
				server.queue.jobs.Push(&work)
			}
		}

		server.queue.DistributeWork()
	}()

	reply = &pb.SimulationReply{
		SimulationId: req.SimulationId,
	}

	return
}
