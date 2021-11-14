package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"math/rand"
	"sync"
)

type taskQueue struct {
	mu       *sync.RWMutex
	newTasks *sync.Cond
	size     uint32
	tasks    []*pb.SimulationRun
	kill     map[uint32]chan bool
	closed   bool
}

func newQueue() (que *taskQueue) {
	mu := &sync.RWMutex{}
	return &taskQueue{
		mu:       mu,
		newTasks: sync.NewCond(mu),
		kill:     make(map[uint32]chan bool),
	}
}

func (que *taskQueue) add(items ...*pb.SimulationRun) {
	que.newTasks.L.Lock()
	defer que.newTasks.L.Unlock()

	que.tasks = append(que.tasks, items...)
	que.size = uint32(len(que.tasks))
	que.newTasks.Broadcast()
}

func (que *taskQueue) pop() (item *pb.SimulationRun, ok bool) {
	que.mu.Lock()
	defer que.mu.Unlock()

	if len(que.tasks) == 0 {
		return
	}

	ok = true
	item, que.tasks = que.tasks[0], que.tasks[1:]
	que.size = uint32(len(que.tasks))

	return
}

func (que *taskQueue) len() (size uint32) {
	que.mu.Lock()
	size = que.size
	que.mu.Unlock()

	return
}

func (que *taskQueue) close() {
	que.mu.Lock()
	defer que.mu.Unlock()

	que.closed = true

	for id, ch := range que.kill {
		ch <- true
		delete(que.kill, id)
	}
}

func (que *taskQueue) linger() (linger bool) {

	que.mu.Lock()
	if que.closed {
		return false
	}
	que.mu.Unlock()

	id := rand.Uint32()

	die := make(chan bool)
	defer close(die)

	que.mu.Lock()
	que.kill[id] = die
	que.mu.Unlock()

	defer func() {
		que.mu.Lock()
		delete(que.kill, id)
		que.mu.Unlock()
	}()

	live := make(chan bool, 1)
	defer close(live)

	go func() {
		que.newTasks.L.Lock()
		que.newTasks.Wait()
		live <- true
		que.newTasks.L.Unlock()
	}()

	select {
	case <-live:
		return true
	case <-die:
		return false
	}
}
