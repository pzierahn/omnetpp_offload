package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (cons *consumer) init(conn *connection) (err error) {

	err = conn.checkout(cons.simulation, cons.simulationTgz)
	if err != nil {
		return
	}

	err = conn.setup(cons.simulation)
	if err != nil {
		return
	}

	stream, err := conn.provider.Allocate(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("[%s] startAllocator", conn.name())

	go func() {

		for {
			alloc, err := stream.Recv()
			if err != nil {
				break
			}

			log.Printf("[%s] allocated %d slots",
				conn.name(), alloc.Slots)

			for inx := uint32(0); inx < alloc.Slots; inx++ {

				task, ok := <-cons.allocator
				if !ok {
					//
					// No tasks left
					//

					return
				}

				go func() {
					// TODO: Find a better way to handle this
					defer cons.finished.Done()

					err := conn.run(task)
					if err != nil {
						log.Fatalln(conn.name(), err)
					}
				}()
			}
		}
	}()

	go func() {
		//
		// Communicate changes in the allocSlots number to the provider
		//

		cond := cons.allocCond

		for {
			// TODO: Change wait and lock
			cond.L.Lock()
			allocateJobs := uint32(len(cons.allocate))
			cond.L.Unlock()

			log.Printf("[%s] request %d slots", conn.name(), allocateJobs)

			err := stream.Send(&pb.AllocateRequest{
				ConsumerId: cons.consumerId,
				Request:    allocateJobs,
			})
			if err != nil {
				log.Println(err)
				break
			}

			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
		}
	}()

	return
}
