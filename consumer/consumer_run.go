package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"sync"
	"time"
)

var taskWG sync.WaitGroup

func (conn *connection) run(task *pb.Simulation) (err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", conn.name(), runName)

	startExec := time.Now()

	resultRef, err := conn.provider.Run(context.Background(), task)
	if err != nil {
		return
	}

	endExec := time.Now()

	log.Printf("[%s] %s finished (%v)",
		conn.name(), runName, endExec.Sub(startExec))

	store := storage.FromClient(conn.store)
	buf, err := store.Download(resultRef)
	if err != nil {
		return
	}

	log.Printf("[%s] %s downloaded results %v in %v",
		conn.name(), runName, simple.ByteSize(uint64(buf.Len())), time.Now().Sub(endExec))

	return
}

func (conn *connection) startScheduler(count <-chan uint32, tasks chan *pb.Simulation) {

	log.Printf("startScheduler: %s", conn.name())

	stream, err := conn.provider.Schedule(context.Background())
	if err != nil {
		return
	}

	go func() {
		for req := range count {
			log.Printf("[%s] request %d", conn.name(), req)

			err := stream.Send(&pb.AllocateRequest{
				Request: req,
			})
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()

	go func() {
		for {
			alloc, err := stream.Recv()
			if err != nil {
				return
			}

			log.Printf("[%s] allocated %d slots", conn.name(), alloc.Slots)

			for inx := uint32(0); inx < alloc.Slots; inx++ {
				task, ok := <-tasks
				if !ok {
					//
					// No tasks left: quit this thread
					//

					return
				}

				go func() {
					// TODO: Find a better way to handle this
					defer taskWG.Done()

					err := conn.run(task)
					if err != nil {
						//
						// Job failed: Reschedule task
						//
						log.Printf("[%s] job %s-%s failed: %v",
							conn.name(), task.Config, task.RunNum, err)

						tasks <- task
					}
				}()
			}
		}
	}()
}

func (cons *consumer) dispatchTasks() (err error) {

	log.Printf("dispatchTasks:")

	var runs *pb.SimulationRuns

	cons.connMu.RLock()
	log.Printf("dispatchTasks: connections=%d", len(cons.connections))

	for _, conn := range cons.connections {
		runs, err = conn.provider.ListRunNums(context.Background(), &pb.Simulation{
			Id:        cons.simulation.Id,
			OppConfig: cons.simulation.OppConfig,

			// TODO: Fix 0 index
			Config: cons.config.SimulateConfigs[0],
		})
		break
	}
	cons.connMu.RUnlock()

	if err != nil {
		return
	}

	tasks := make(chan *pb.Simulation)
	defer close(tasks)

	counters := make([]chan uint32, 0)
	defer func() {
		for _, counter := range counters {
			close(counter)
		}
	}()

	for _, conn := range cons.connections {
		counter := make(chan uint32)
		counters = append(counters, counter)
		go conn.startScheduler(counter, tasks)
	}

	log.Printf("dispatchTasks: counters=%d", len(runs.Runs))
	for _, counter := range counters {
		counter <- uint32(len(runs.Runs))
	}

	for inx, num := range runs.Runs {
		log.Printf("dispatchTasks: num=%s", num)
		taskWG.Add(1)
		tasks <- &pb.Simulation{
			Id:        cons.simulation.Id,
			OppConfig: cons.simulation.OppConfig,
			Config:    runs.Config,
			RunNum:    num,
		}

		for _, counter := range counters {
			counter <- uint32(len(runs.Runs) - inx - 1)
		}
	}

	taskWG.Wait()

	return
}
