package broker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) SimNew(_ context.Context, simulation *pb.Simulation) (resp *pb.SimulationId, err error) {

	sState := server.simulations.createNew(simulation)
	logger.Printf("created new simulation: id='%s'", sState.simulationId)

	resp = &pb.SimulationId{
		Id: sState.simulationId,
	}

	return
}

func (server *broker) SimAddTasks(_ context.Context, tasks *pb.Tasks) (resp *pb.Empty, err error) {

	logger.Printf("simulation %s (added %d tasks)", tasks.SimulationId, len(tasks.Items))

	sState := server.simulations.getSimulationState(tasks.SimulationId)
	sState.addTasks(tasks.Items...)

	go server.distribute()

	resp = &pb.Empty{}

	return
}

func (server *broker) SetSource(_ context.Context, ref *pb.Source) (resp *pb.Empty, err error) {

	logger.Printf("set source for %s to %v", ref.SimulationId, ref.Source)

	sState := server.simulations.getSimulationState(ref.SimulationId)
	sState.setSource(ref.Source)

	resp = &pb.Empty{}

	return
}

func (server *broker) GetSource(_ context.Context, sim *pb.SimulationId) (resp *pb.StorageRef, err error) {

	logger.Printf("get source for %s", sim.Id)

	sState := server.simulations.getSimulationState(sim.Id)
	resp = sState.getSource()

	return
}

func (server *broker) AddBinary(_ context.Context, binary *pb.SimBinary) (resp *pb.Empty, err error) {

	logger.Printf("%s: new binary (%s)", binary.SimulationId, binary.Arch)

	sState := server.simulations.getSimulationState(binary.SimulationId)
	sState.addBinary(binary)

	resp = &pb.Empty{}

	return
}
