package provider

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (client *workerConnection) compile(simulation *pb.Simulation) {
	simulationId := simulation.SimulationId

	simulationBase, err := client.checkout(simulationId)
	if err != nil {
		panic(err)
	}

	compiler := Compiler{
		Broker:         client.broker,
		Storage:        client.storage,
		SimulationId:   simulation.SimulationId,
		SimulationBase: simulationBase,
		OppConfig:      simulation.OppConfig,
	}

	err = compiler.Compile()
	if err != nil {
		panic(err)
	}

	err = compiler.CheckinBinary()
	if err != nil {
		panic(err)
	}
}
