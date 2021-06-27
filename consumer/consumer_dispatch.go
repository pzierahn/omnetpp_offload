package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (cons *consumer) dispatchTasks() (err error) {

	cond := cons.allocCond

	// Wait for initialisation of the queue
	cons.allocCond.L.Lock()
	cons.allocCond.Wait()
	cons.allocCond.L.Unlock()

	log.Printf("dispatchTasks:")

	for {
		var schedule *pb.SimulationRun

		cond.L.Lock()

		log.Printf("dispatchTasks: left tasks %d", len(cons.allocate))

		if len(cons.allocate) == 0 {
			cond.L.Unlock()
			break
		}

		schedule, cons.allocate = cons.allocate[0], cons.allocate[1:]
		cond.Broadcast()
		cond.L.Unlock()

		log.Printf("dispatchTasks: num=%s", schedule.RunNum)
		cons.finished.Add(1)

		cons.allocator <- schedule
	}

	cons.finished.Wait()

	return
}
