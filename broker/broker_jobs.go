package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"sync"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    *pb.Work
	priority int
	index    int
}

// An IntHeap is a min-heap of ints.
type WorkHeap []*pb.Work

func (h WorkHeap) Len() int           { return len(h) }
func (h WorkHeap) Less(i, j int) bool { return h[i].SimulationId < h[j].SimulationId }
func (h WorkHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WorkHeap) Push(x *pb.Work) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x)
}

func (h *WorkHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type queue struct {
	mu      sync.Mutex
	jobs    WorkHeap
	workers map[string]chan *pb.Work
}

func (que *queue) Link(id string, worker chan *pb.Work) {
	que.mu.Lock()

	logger.Println("link", id)
	que.workers[id] = worker

	que.mu.Unlock()
}

func (que *queue) Unlink(id string) {
	que.mu.Lock()

	logger.Println("unlink", id)
	delete(que.workers, id)

	que.mu.Unlock()
}

func (que *queue) DistributeWork() {
	que.mu.Lock()

	for que.jobs.Len() > 0 {
		job := que.jobs.Pop()
		work := job.(*pb.Work)

		for workerId, stream := range que.workers {
			logger.Println("assign", work.ConfigId, "-->", workerId)
			stream <- work
			break
		}
	}

	que.mu.Unlock()
}
