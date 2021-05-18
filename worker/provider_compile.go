package worker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (client *workerConnection) compile(compileJob *pb.Work_Compile) {
	simulation := compileJob.Compile
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
