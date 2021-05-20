package provider

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (client *workerConnection) compile(build *pb.Assignment_Build) {
	simulationId := build.Build.SimulationId

	simulationBase, err := client.checkout(simulationId)
	if err != nil {
		panic(err)
	}

	compiler := Compiler{
		Broker:         client.broker,
		Storage:        client.storage,
		SimulationId:   simulationId,
		SimulationBase: simulationBase,
		OppConfig:      build.Build.OppConfig,
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
