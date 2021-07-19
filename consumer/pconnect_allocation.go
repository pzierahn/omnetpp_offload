package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (pConn *providerConnection) allocationHandler(stream pb.Provider_AllocateClient, cons *consumer) {
	log.Printf("[%s] start allocator", pConn.name())

	for {
		alloc, err := stream.Recv()
		if err != nil {
			break
		}

		log.Printf("[%s] allocated %d slots",
			pConn.name(), alloc.Slots)

		for inx := uint32(0); inx < alloc.Slots; inx++ {

			task, ok := cons.allocate.pop()
			if !ok {
				//
				// No tasks left
				//

				return
			}

			go func() {
				// TODO: Find a better way to handle this

				err := pConn.run(task)
				if err != nil {
					log.Printf("[%s] error %v", pConn.name(), err)
					log.Printf("[%s] reschedule %s_%s", pConn.name(), task.Config, task.RunNum)

					// Add item back to queue to send right allocation num
					cons.allocate.add(task)
				} else {
					cons.finished.Done()
				}
			}()
		}
	}
}

func (pConn *providerConnection) sendAllocationRequest(stream pb.Provider_AllocateClient, cons *consumer) (err error) {
	request := cons.allocate.len()
	log.Printf("[%s] request %d slots", pConn.name(), request)

	err = stream.Send(&pb.AllocateRequest{
		SimulationId: cons.simulation.Id,
		Request:      uint32(request),
	})

	return
}
