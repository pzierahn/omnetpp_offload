package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"sync"
)

type simulationManager struct {
	sync.RWMutex
	simulations map[string]*simulationState
}

func newSimulationManager() (state simulationManager) {
	return simulationManager{
		simulations: make(map[string]*simulationState),
	}
}

func (sm *simulationManager) createNew(simulation *pb.Simulation) (sState *simulationState) {

	sState = newSimulationState(simulation)

	sm.Lock()
	sm.simulations[sState.simulationId] = sState
	sm.Unlock()

	return
}

func (sm *simulationManager) getSimulationState(id string) (sState *simulationState) {

	sm.RLock()
	defer sm.RUnlock()

	sState = sm.simulations[id]

	return
}

func (sm *simulationManager) pullCompile(arch *pb.OsArch) (simulation *pb.Simulation) {

	sm.RLock()
	defer sm.RUnlock()

	for _, sim := range sm.simulations {
		_, ok := sim.binaries[osArchId(arch)]

		if !ok {

			//
			// Binary is not compiled for arch
			//

			simulation = &pb.Simulation{
				SimulationId: sim.simulationId,
				OppConfig:    sim.oppConfig,
			}
			break
		}
	}

	return
}

func (sm *simulationManager) pullWork() (task *pb.Task) {

	sm.RLock()
	defer sm.RUnlock()

	for _, sim := range sm.simulations {
		var id taskId
		for id, task = range sim.tasks {
			// TODO: Fix this mess
			delete(sim.tasks, id)
			return
		}
	}

	return
}
