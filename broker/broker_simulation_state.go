package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"strings"
	"sync"
)

type osArch string

func osArchId(binary *pb.OsArch) osArch {
	id := strings.Join([]string{
		binary.Os,
		binary.Arch,
	}, "_")

	return osArch(id)
}

type taskId string

func tId(task *pb.Task) taskId {
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
	tasks        map[taskId]*pb.Task
	// assignments  map[workerId]*pb.Task
	finished  map[taskId]*pb.Task
	source    *pb.StorageRef
	oppConfig *pb.OppConfig
	binaries  map[osArch]*pb.SimBinary
}

func newSimulationState(simulation *pb.Simulation) (state *simulationState) {
	return &simulationState{
		simulationId: simple.NamedId(simulation.Tag, 8),
		tasks:        make(map[taskId]*pb.Task),
		finished:     make(map[taskId]*pb.Task),
		oppConfig:    simulation.OppConfig,
		binaries:     make(map[osArch]*pb.SimBinary),
	}
}

func (ss *simulationState) addTasks(tasks ...*pb.Task) {

	ss.Lock()
	defer ss.Unlock()

	for _, task := range tasks {
		ss.tasks[tId(task)] = task
	}

	return
}

func (ss *simulationState) setSource(ref *pb.StorageRef) {

	ss.Lock()
	defer ss.Unlock()

	ss.source = ref

	return
}

func (ss *simulationState) getSource() (ref *pb.StorageRef) {

	ss.RLock()
	defer ss.RUnlock()

	ref = ss.source

	return
}

func (ss *simulationState) addBinary(binary *pb.SimBinary) {

	ss.Lock()
	defer ss.Unlock()

	ss.binaries[osArchId(binary.Arch)] = binary

	return
}
