package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
)

func (server *broker) distributeWork() {
	server.queue.mu.Lock()
	server.workers.Lock()

	logger.Printf("distributeWork jobs %d\n", server.queue.jobs.Len())
	logger.Printf("available workers %d\n", len(server.queue.workers))

	for workerId, stream := range server.queue.workers {

		status, _ := server.workers.workers[workerId]

		logger.Printf("Checking %s --> %d\n", workerId, status.FreeResources)

		if status.FreeResources == 0 {
			//
			// Client is busy
			//

			logger.Printf("Client busy %s\n", workerId)

			continue
		}

		packages := simple.MathMin(
			server.queue.jobs.Len(),
			int(status.FreeResources),
		)

		for inx := 0; inx < packages; inx++ {
			job := server.queue.jobs.Pop()
			work := job.(*pb.Work)

			logger.Printf("assign %s.%s --> %s\n",
				work.SimulationId, work.ConfigId, workerId)
			stream <- work
		}
	}

	server.workers.Unlock()
	server.queue.mu.Unlock()
}
