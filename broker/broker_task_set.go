package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"strings"
)

func taskId(task *pb.Task) (id string) {

	id = strings.Join([]string{
		task.SimulationId,
		task.Config,
		task.RunNumber,
	}, ".")

	return
}

type taskSet struct {
	list map[string]*pb.Task
}

func (set *taskSet) pop(length int) (tasks []*pb.Task) {
	for key, elem := range set.list {

		if length == 0 {
			break
		}

		tasks = append(tasks, elem)
		delete(set.list, key)

		length--
	}

	return
}

func (set *taskSet) len() (length int) {
	length = len(set.list)
	return
}

func (set *taskSet) merge(set2 *taskSet) {

	if set.list == nil {
		set.list = make(map[string]*pb.Task)
	}

	set.add(set2.pop(-1)...)
}

func (set *taskSet) add(tasks ...*pb.Task) {

	if set.list == nil {
		set.list = make(map[string]*pb.Task)
	}

	for _, task := range tasks {
		set.list[taskId(task)] = task
	}
}

func (set *taskSet) remove(tasks ...*pb.Task) {
	for _, task := range tasks {
		delete(set.list, taskId(task))
	}
}

type simulationState struct {
	queue    taskSet
	assigned map[string]*taskSet // workerId --> tasks
	finished []*pb.TaskResult
}

func newSimulationState() (state simulationState) {

	state = simulationState{
		assigned: make(map[string]*taskSet),
	}

	return
}

func (state *simulationState) assign(workerId string, tasks ...*pb.Task) {

	assignments, ok := state.assigned[workerId]

	if !ok {
		assignments = &taskSet{}
	}

	assignments.add(tasks...)

	state.assigned[workerId] = assignments

	return
}

func (state *simulationState) finish(result *pb.TaskResult) {
	state.finished = append(state.finished, result)

	for _, set := range state.assigned {
		set.remove(result.Task)
	}

	return
}
