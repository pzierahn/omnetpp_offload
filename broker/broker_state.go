package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"sync"
)

type taskState struct {
	sync.RWMutex
	assignments map[string][]*pb.Task
	capacities  map[string]*pb.ResourceCapacity
	tasks       WorkHeap
	workers     map[string]chan<- *pb.Tasks
}

func initTasksDB() (state taskState) {
	state = taskState{
		assignments: make(map[string][]*pb.Task),
		capacities:  make(map[string]*pb.ResourceCapacity),
		tasks:       WorkHeap{},
		workers:     make(map[string]chan<- *pb.Tasks),
	}

	return
}

func (state *taskState) SetCapacity(id string, cap *pb.ResourceCapacity) {
	state.Lock()
	defer state.Unlock()

	state.capacities[id] = cap

	return
}

func (state *taskState) NewWorker(id string) (worker chan *pb.Tasks) {
	state.Lock()
	defer state.Unlock()

	logger.Printf("new worker %v\n", id)
	worker = make(chan *pb.Tasks)
	state.workers[id] = worker

	return
}

func (state *taskState) RemoveWorker(id string) {
	state.Lock()
	defer state.Unlock()

	logger.Printf("remove worker %v\n", id)

	if ch, ok := state.workers[id]; ok && ch != nil {
		close(ch)
	}
	delete(state.workers, id)

	if tasks, ok := state.assignments[id]; ok && len(tasks) > 0 {
		//
		// Connection lost without finishing all assigned tasks!
		//

		logger.Printf("reassign %d unfinished jobs from %s\n", len(tasks), id)

		state.tasks.Push(tasks...)
	}

	delete(state.assignments, id)
	delete(state.capacities, id)

	return
}

func (state *taskState) AddTasks(tasks ...*pb.Task) {
	state.Lock()
	defer state.Unlock()

	state.tasks.Push(tasks...)

	return
}

func (state *taskState) DistributeWork() {
	state.Lock()
	defer state.Unlock()

	logger.Printf("distribute work (%d workers, %d tasks)\n",
		len(state.workers), state.tasks.Len())

	if state.tasks.Len() == 0 {
		return
	}

	for workerId, stream := range state.workers {

		capacity, ok := state.capacities[workerId]

		if !ok || capacity.FreeResources <= 0 {
			//
			// Client is busy
			//

			logger.Printf("%s busy\n", workerId)

			continue
		}

		logger.Printf("%s capacity %d\n", workerId, capacity.FreeResources)

		packages := simple.MathMin(
			state.tasks.Len(),
			int(capacity.FreeResources),
		)

		var jobs []*pb.Task

		for inx := 0; inx < packages; inx++ {
			task := state.tasks.Pop()
			jobs = append(jobs, task)
		}

		logger.Printf("assign %v --> %v\n", workerId, packages)

		state.assignments[workerId] = jobs

		// Send data to worker
		stream <- &pb.Tasks{Items: jobs}

		logger.Printf("assign done %v --> %v\n", workerId, packages)

		// Remove client info from worker queue
		delete(state.capacities, workerId)
	}
}
