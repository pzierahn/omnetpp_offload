package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"sync"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    *pb.Task
	priority int
	index    int
}

// An IntHeap is a min-heap of ints.
type WorkHeap []*pb.Task

func (h WorkHeap) Len() int {
	return len(h)
}

func (h WorkHeap) Less(i, j int) bool {
	return h[i].Config+h[i].RunNumber < h[j].Config+h[j].RunNumber
}

func (h WorkHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *WorkHeap) Push(x *pb.Task) {
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
	workers map[string]chan<- *pb.Tasks
}

func (que *queue) Link(id string, worker chan<- *pb.Tasks) {
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

func initQueue() (que queue) {
	que = queue{
		mu:      sync.Mutex{},
		jobs:    make(WorkHeap, 0),
		workers: make(map[string]chan<- *pb.Tasks),
	}

	return
}
