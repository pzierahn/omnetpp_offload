package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"sync"
)

type distributor struct {
	sync.RWMutex
	capacities  map[string]*pb.ResourceCapacity // workerId --> ResourceCapacity
	workers     map[string]chan<- *pb.Task      // workerId --> task channel
	simulations map[string]*simulationState     // simulationId --> state
}

func initTasksDB() (state distributor) {
	state = distributor{
		capacities:  make(map[string]*pb.ResourceCapacity),
		workers:     make(map[string]chan<- *pb.Task),
		simulations: make(map[string]*simulationState),
	}

	return
}

func (state *distributor) ReceiveResult(result *pb.TaskResult) {
	state.Lock()
	defer state.Unlock()

	sim := state.simulations[result.Task.SimulationId]
	sim.finish(result)
}

func (state *distributor) SetCapacity(id string, cap *pb.ResourceCapacity) {
	state.Lock()
	defer state.Unlock()

	state.capacities[id] = cap

	return
}

func (state *distributor) NewWorker(id string) (worker chan *pb.Task) {
	state.Lock()
	defer state.Unlock()

	logger.Printf("new worker %v\n", id)
	worker = make(chan *pb.Task)
	state.workers[id] = worker

	return
}

func (state *distributor) RemoveWorker(workerId string) {
	state.Lock()
	defer state.Unlock()

	logger.Printf("remove worker %v\n", workerId)

	if ch, ok := state.workers[workerId]; ok && ch != nil {
		close(ch)
	}
	delete(state.workers, workerId)

	for id, simulation := range state.simulations {

		logger.Printf("checking for unfinished jobs for %v in %v\n", workerId, id)

		if unFinished, ok := simulation.assigned[workerId]; ok && unFinished.len() > 0 {
			//
			// Connection lost without finishing all assigned tasks!
			//

			logger.Printf("reassign %d unfinished jobs from %s\n", unFinished.len(), workerId)

			simulation.queue.merge(unFinished)
		}

		delete(simulation.assigned, workerId)
	}

	delete(state.capacities, workerId)

	return
}

func (state *distributor) NewSimulation(req *pb.Simulation) {

	simulation := newSimulationState()

	for _, config := range req.Run {
		for _, runNum := range config.RunNumbers {
			task := &pb.Task{
				SimulationId: req.SimulationId,
				OppConfig:    req.OppConfig,
				Source:       req.Source,
				Config:       config.Config,
				RunNumber:    runNum,
			}

			simulation.queue.add(task)
		}
	}

	state.Lock()
	state.simulations[req.SimulationId] = &simulation
	state.Unlock()

	return
}

func (state *distributor) DistributeWork() {
	state.Lock()
	defer state.Unlock()

	simulations := 0
	jobCount := 0

	for _, sim := range state.simulations {
		simulations++
		jobCount += sim.queue.len()
	}

	logger.Printf("distribute work (%d workers, %d simulations, %d jobs)\n",
		len(state.workers), simulations, jobCount)

	for workerId, stream := range state.workers {

		capacity, ok := state.capacities[workerId]

		if !ok || capacity.FreeResources <= 0 {
			//
			// Client is busy
			//

			logger.Printf("%s busy\n", workerId)

			continue
		}

		var simulation *simulationState

		for _, sim := range state.simulations {
			simulation = sim

			if sim.queue.len() > 0 {
				break
			}
		}

		if simulation == nil {

			//
			// No simulations left
			//

			break
		}

		logger.Printf("%s capacity %d\n", workerId, capacity.FreeResources)

		packages := simple.MathMin(
			simulation.queue.len(),
			int(capacity.FreeResources),
		)

		tasks := simulation.queue.pop(packages)

		for _, task := range tasks {
			stream <- task
		}

		simulation.assign(workerId, tasks...)

		//tasks := simulation.queue.pop(packages)
		//simulation.assign(workerId, tasks...)
		//
		//logger.Printf("sending %v %v tasks\n", workerId, len(tasks))
		//
		//// Send data to worker
		//stream <- &pb.Tasks{Items: tasks}

		// Remove client info from worker queue
		delete(state.capacities, workerId)
	}
}

func (state *distributor) Status(req *pb.ResultsRequest) (reply *pb.StatusReply, err error) {
	state.RLock()
	defer state.RUnlock()

	reply = &pb.StatusReply{
		SimulationId: req.SimulationId,
	}

	sim, ok := state.simulations[req.SimulationId]

	if !ok {
		return
	}

	for _, elem := range sim.queue.list {
		reply.Queue = append(reply.Queue, &pb.StatusReply_QueueInfo{
			Config: elem.Config,
			RunNum: elem.RunNumber,
		})
	}

	for workerId, ass := range sim.assigned {
		for _, elem := range ass.list {
			info := &pb.StatusReply_QueueInfo{
				Config: elem.Config,
				RunNum: elem.RunNumber,
			}

			reply.Assigned = append(reply.Assigned, &pb.StatusReply_Assignment{
				Config:   info,
				WorkerId: workerId,
			})
		}
	}

	for _, elem := range sim.finished {
		reply.Finished = append(reply.Finished, &pb.StatusReply_QueueInfo{
			Config: elem.Task.Config,
			RunNum: elem.Task.RunNumber,
		})
	}

	return
}
