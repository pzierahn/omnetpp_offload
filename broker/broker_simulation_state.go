package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"strings"
	"sync"
)

type osArch string

func osArchId(binary *pb.Arch) osArch {
	id := strings.Join([]string{
		binary.Os,
		binary.Arch,
	}, "_")

	return osArch(id)
}

type taskId string

func tId(task *pb.SimulationRun) taskId {
	id := strings.Join([]string{
		task.SimulationId,
		task.Config,
		task.RunNumber,
	}, "_")

	return taskId(id)
}

type simulationState struct {
	sync.RWMutex
	simulationId string
	tasks        map[taskId]*pb.SimulationRun
	finished     map[taskId]*pb.SimulationRun
	source       *pb.StorageRef
	oppConfig    *pb.OppConfig
	binaries     map[osArch]*pb.Binary
	// assignments  map[workerId]*pb.Task
}

func newSimulationState(simulation *pb.Simulation) (state *simulationState) {
	return &simulationState{
		simulationId: simple.NamedId(simulation.Tag, 8),
		tasks:        make(map[taskId]*pb.SimulationRun),
		finished:     make(map[taskId]*pb.SimulationRun),
		oppConfig:    simulation.OppConfig,
		binaries:     make(map[osArch]*pb.Binary),
	}
}

func (ss *simulationState) write(fn func()) {

	ss.Lock()
	defer ss.Unlock()

	fn()

	return
}

func (ss *simulationState) read(fn func()) {

	ss.RLock()
	defer ss.RUnlock()

	fn()

	return
}
