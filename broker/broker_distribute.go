package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
)

func (server *broker) distributeWork() {
	server.queue.mu.Lock()
	server.workers.Lock()

	defer func() {
		server.workers.Unlock()
		server.queue.mu.Unlock()
	}()

	logger.Printf("jobs %d in queue\n", server.queue.jobs.Len())
	logger.Printf("available workers %d\n", len(server.queue.workers))

	if server.queue.jobs.Len() == 0 {
		return
	}

	for workerId, stream := range server.queue.workers {

		status, ok := server.workers.workers[workerId]

		if !ok {
			logger.Printf("Checking %s --> busy\n", workerId)
			continue
		}

		logger.Printf("Checking %s --> %d\n", workerId, status.FreeResources)

		if status.FreeResources <= 0 {
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

		var jobs []*pb.Work

		for inx := 0; inx < packages; inx++ {
			job := server.queue.jobs.Pop()
			work := job.(*pb.Work)

			jobs = append(jobs, work)
		}

		logger.Printf("assign %s --> %s\n", workerId, jobs)

		// Send data to worker
		stream <- &pb.Tasks{Jobs: jobs}

		// Remove client info from worker queue
		delete(server.workers.workers, workerId)
	}
}
