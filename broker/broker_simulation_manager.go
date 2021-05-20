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

func (sm *simulationManager) pullCompile(arch *pb.Arch) (sState *simulationState) {

	sm.RLock()
	defer sm.RUnlock()

	for _, sim := range sm.simulations {

		if sim.source == nil {
			// No simulation source: skip
			continue
		}

		_, ok := sim.binaries[osArchId(arch)]

		if !ok {

			//
			// Simulation is not compiled for arch
			//

			sState = sim

			break
		}
	}

	return
}

func (sm *simulationManager) pullWork(arch *pb.Arch) (task *pb.SimulationRun) {

	sm.RLock()
	defer sm.RUnlock()

	for _, sim := range sm.simulations {
		if sim.source == nil {
			// No simulation source: skip
			continue
		}

		_, ok := sim.binaries[osArchId(arch)]

		if !ok {
			// No simulation binary: skip
			continue
		}

		var id taskId
		for id = range sim.queue {
			delete(sim.queue, id)
			task = sim.runs[id]

			return
		}
	}

	return
}
