package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/provider"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (connect *providerConnection) execute(sim *simulation) (err error) {

	downloadQueue := make(chan *download, 32)
	defer close(downloadQueue)

	go connect.resultsDownloader(downloadQueue, sim)

	ctx := metadata.AppendToOutgoingContext(sim.ctx, provider.MetaSimulationId, sim.id)

	stream, err := connect.provider.Allocate(ctx)
	if err != nil {
		return
	}

	log.Printf("[%s] start allocator", connect.id())

	for {
		_, err = stream.Recv()
		if err != nil {
			log.Printf("[%s] error: %v", connect.id(), err)
			break
		}

		log.Printf("[%s] allocated slot", connect.id())

		task, ok := sim.queue.pop()
		if !ok {
			//
			// No tasks left in queue.
			//

			_ = stream.Send(&pb.FreeSlot{})

			log.Printf("[%s] idle", connect.id())

			// Wait to see if more slots are needed.
			if sim.queue.linger() {

				//
				// A task was rescheduled, continue requesting slots.
				//

				log.Printf("[%s] continue requesting slots", connect.id())

				continue
			} else {

				//
				// Simulation is finished, stop requesting slots.
				//

				log.Printf("[%s] stop requesting slots", connect.id())

				_ = stream.CloseSend()
				_, _ = connect.provider.StopEvaluation(ctx, &emptypb.Empty{})

				break
			}
		}

		go func() {
			ref, err := connect.run(task)
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

				log.Printf("[%s] run failed: reschedule %+v", connect.id(), task)
				sim.queue.add(task)
			}
		}()
	}

	return
}
