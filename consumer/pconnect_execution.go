package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (pConn *providerConnection) allocationHandler(stream pb.Provider_AllocateClient, sim *simulation) {
	log.Printf("[%s] start allocator", pConn.id())

	for {
		alloc, err := stream.Recv()
		if err != nil {
			break
		}

		log.Printf("[%s] allocated %d slots",
			pConn.id(), alloc.Slots)

		for inx := uint32(0); inx < alloc.Slots; inx++ {

			task, ok := sim.queue.pop()
			if !ok {
				//
				// No tasks left
				//

				return
			}

			go func() {
				ref, err := pConn.run(task)
				if err != nil {
					log.Printf("[%s] run failed: reschedule %+v", pConn.id(), task)
					sim.queue.add(task)
					return
				}

				pConn.downloadPipe <- &download{
					task: task,
					ref:  ref,
				}

				//buf, err := pConn.download(ref)
				//if err != nil {
				//	log.Printf("[%s] download failed: reschedule %+v", pConn.id(), task)
				//	// Add item back to taskQueue to send right allocation num
				//	cons.queue.add(task)
				//	return
				//}
				//
				//done := eval.LogAction(eval.ActionExtract, ref.Filename)
				//cons.extractResults(buf)
				//_ = done(nil)
				//
				//_, _ = pConn.store.Delete(pConn.ctx, ref)
				//
				//cons.finished.Done()
			}()
		}
	}
}

func (pConn *providerConnection) execute(sim *simulation) (err error) {
	go pConn.downloader(1, sim)

	stream, err := pConn.provider.Allocate(sim.ctx)
	if err != nil {
		return
	}

	go pConn.allocationHandler(stream, sim)

	go sim.queue.onChange(func() (cancel bool) {
		request := sim.queue.len()
		log.Printf("[%s] request %d slots", pConn.id(), request)

		err = stream.Send(&pb.AllocateRequest{
			SimulationId: sim.id,
			Request:      uint32(request),
		})

		if err != nil {
			log.Println(err)
			return true
		}

		return false
	})

	return
}