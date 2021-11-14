package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/provider"
	"google.golang.org/grpc/metadata"
	"log"
)

func (pConn *providerConnection) execute(sim *simulation) (err error) {

	downloadQueue := make(chan *download, 32)
	defer close(downloadQueue)

	go pConn.resultsDownloader(downloadQueue, sim)

	ctx := metadata.AppendToOutgoingContext(sim.ctx, provider.MetaSimulationId, sim.id)

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

			log.Printf("[%s] idle", pConn.id())

			// Wait to see if more slots are needed.
			if sim.queue.linger() {

				//
				// A task was rescheduled, continue requesting slots.
				//

				log.Printf("[%s] continue requesting slots", pConn.id())

				continue
			} else {

				//
				// Simulation is finished, stop requesting slots.
				//

				log.Printf("[%s] stop requesting slots", pConn.id())

				_ = stream.CloseSend()

				break
			}
		}

		go func() {
			ref, err := pConn.run(task)
			_ = stream.Send(&pb.FreeSlot{})

			if err == nil {

				//
				// Execution was a success: download result
				//

				downloadQueue <- &download{
					task: task,
					ref:  ref,
				}
			} else {

				//
				// Execution failed: reschedule
				//

				log.Printf("[%s] run failed: reschedule %+v", pConn.id(), task)
				sim.queue.add(task)
			}
		}()
	}

	return
}
