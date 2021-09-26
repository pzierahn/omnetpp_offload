package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/grpc/metadata"
	"log"
)

func (pConn *providerConnection) execute(sim *simulation) (err error) {
	go pConn.resultsDownloader(sim)

	ctx := metadata.AppendToOutgoingContext(sim.ctx, "simulationId", sim.id)

	stream, err := pConn.provider.Allocate(ctx)
	if err != nil {
		return
	}

	log.Printf("[%s] start allocator", pConn.id())

	for {
		_, err = stream.Recv()
		if err != nil {
			log.Printf("[%s] error: %v", pConn.id(), err)
			break
		}

		log.Printf("[%s] allocated slot", pConn.id())

		task, ok := sim.queue.pop()
		if !ok {
			//
			// No tasks left in queue.
			//

			_ = stream.Send(&pb.FreeSlot{})

			// Wait if more slots are required.
			if sim.queue.linger() {

				//
				// A task was rescheduled, continue requesting slots.
				//

				continue
			} else {

				//
				// Simulation is finished, stop requesting slots.
				//

				log.Printf("[%s] stop requesting slots", pConn.id())

				break
			}
		}

		go func() {
			ref, err := pConn.run(task)
			_ = stream.Send(&pb.FreeSlot{})

			if err != nil {
				log.Printf("[%s] run failed: reschedule %+v", pConn.id(), task)
				sim.queue.add(task)
				return
			}

			pConn.downloadQueue <- &download{
				task: task,
				ref:  ref,
			}
		}()
	}

	return
}
