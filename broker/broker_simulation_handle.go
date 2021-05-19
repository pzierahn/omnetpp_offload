package broker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (server *broker) Create(_ context.Context, simulation *pb.Simulation) (resp *pb.SimulationId, err error) {

	sState := server.simulations.createNew(simulation)
	logger.Printf("created new simulation: id='%s'", sState.simulationId)

	resp = &pb.SimulationId{
		Id: sState.simulationId,
	}

	return
}

func (server *broker) AddTasks(_ context.Context, tasks *pb.Tasks) (resp *pb.Empty, err error) {

	logger.Printf("simulation %s (added %d tasks)", tasks.SimulationId, len(tasks.Items))

	sState := server.simulations.getSimulationState(tasks.SimulationId)
	sState.write(func() {
		for _, task := range tasks.Items {
			id := tId(task)
			sState.queue[id] = true
			sState.runs[id] = task
		}
	})

	resp = &pb.Empty{}

	return
}

func (server *broker) SetSource(_ context.Context, ref *pb.Source) (resp *pb.Empty, err error) {

	logger.Printf("set source for %s to %v", ref.SimulationId, ref.Source)

	sState := server.simulations.getSimulationState(ref.SimulationId)
	sState.write(func() {
		sState.source = ref.Source
	})

	resp = &pb.Empty{}

	return
}

func (server *broker) GetSource(_ context.Context, sim *pb.SimulationId) (resp *pb.Source, err error) {

	logger.Printf("get source for %s", sim.Id)

	ch := make(chan *pb.Source)
	defer close(ch)

	sState := server.simulations.getSimulationState(sim.Id)
	sState.read(func() {
		ch <- &pb.Source{
			SimulationId: sState.simulationId,
			Source:       sState.source,
		}
	})

	resp = <-ch

	return
}

func (server *broker) AddBinary(_ context.Context, binary *pb.Binary) (resp *pb.Empty, err error) {

	logger.Printf("%s: new binary (%s)", binary.SimulationId, binary.Arch)

	sState := server.simulations.getSimulationState(binary.SimulationId)
	sState.write(func() {
		sState.binaries[osArchId(binary.Arch)] = binary
	})

	resp = &pb.Empty{}

	return
}
