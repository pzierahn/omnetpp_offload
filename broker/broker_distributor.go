package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/worker"
	"sync"
)

type distributor struct {
	sync.RWMutex
	capacities  map[string]int               // workerId --> ResourceCapacity
	workerInfo  map[string]worker.DeviceInfo // workerId --> ResourceCapacity
	workers     map[string]chan<- *pb.Task   // workerId --> task channel
	simulations map[string]*simulationState  // simulationId --> state
}

func initTasksDB() (state distributor) {
	state = distributor{
		capacities:  make(map[string]int),
		workerInfo:  make(map[string]worker.DeviceInfo),
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

func (state *distributor) IncreaseCapacity(id string) {

	state.Lock()
	state.capacities[id]++
	state.Unlock()

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

		capacity := state.capacities[workerId]

		if capacity <= 0 {
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

		logger.Printf("%s capacity %d\n", workerId, capacity)

		packages := simple.MathMin(
			simulation.queue.len(),
			capacity,
		)

		tasks := simulation.queue.pop(packages)

		logger.Printf("sending %d tasks to %s\n", packages, workerId)

		for _, task := range tasks {
			stream <- task
		}

		logger.Printf("assign tasks to %s\n", workerId)
		simulation.assign(workerId, tasks...)

		state.capacities[workerId] -= packages
		logger.Printf("set %s capacity to %d", workerId, state.capacities[workerId])

		//logger.Printf("remove capacities reference %s\n", workerId)
		//// Remove client info from worker queue
		//delete(state.capacities, workerId)
	}
}

func (state *distributor) Status(req *pb.StatusRequest) (reply *pb.StatusReply, err error) {
	state.RLock()
	defer state.RUnlock()

	ids := req.SimulationIds

	if len(ids) == 0 {
		for id := range state.simulations {
			ids = append(ids, id)
		}
	}

	reply = &pb.StatusReply{}

	for _, id := range ids {
		sim, ok := state.simulations[id]

		if !ok {
			continue
		}

		var queue []*pb.SimulationStatus_QueueInfo
		var assignments []*pb.SimulationStatus_Assignment
		var finished []*pb.SimulationStatus_QueueInfo

		for _, elem := range sim.queue.list {
			queue = append(queue, &pb.SimulationStatus_QueueInfo{
				Config: elem.Config,
				RunNum: elem.RunNumber,
			})
		}

		for workerId, ass := range sim.assigned {
			for _, elem := range ass.list {
				info := &pb.SimulationStatus_QueueInfo{
					Config: elem.Config,
					RunNum: elem.RunNumber,
				}

				assignments = append(assignments, &pb.SimulationStatus_Assignment{
					Config:   info,
					WorkerId: workerId,
				})
			}
		}

		for _, elem := range sim.finished {
			finished = append(finished, &pb.SimulationStatus_QueueInfo{
				Config: elem.Task.Config,
				RunNum: elem.Task.RunNumber,
			})
		}

		reply.Items = append(reply.Items, &pb.SimulationStatus{
			SimulationId: id,
			Queue:        queue,
			Assigned:     assignments,
			Finished:     finished,
		})
	}

	return
}
