package broker

import (
	"fmt"
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

	id := fmt.Sprintf("%s_%s_%03s",
		task.SimulationId, task.Config, task.RunNumber)

	return taskId(id)
}

type simulationState struct {
	sync.RWMutex
	simulationId string
	queue        map[taskId]bool
	runs         map[taskId]*pb.SimulationRun
	source       *pb.StorageRef
	oppConfig    *pb.OppConfig
	binaries     map[osArch]*pb.Binary
	// assignments  map[workerId]*pb.Task
}

func newSimulationState(simulation *pb.Simulation) (state *simulationState) {
	return &simulationState{
		simulationId: simple.NamedId(simulation.Tag, 8),
		queue:        make(map[taskId]bool),
		runs:         make(map[taskId]*pb.SimulationRun, 0),
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
