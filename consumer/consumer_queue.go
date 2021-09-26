package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"math/rand"
	"sync"
)

type taskQueue struct {
	mu    *sync.RWMutex
	cond  *sync.Cond
	size  uint32
	tasks []*pb.SimulationRun
	kill  map[uint32]chan bool
}

func newQueue() (que *taskQueue) {
	mu := &sync.RWMutex{}
	return &taskQueue{
		mu:   mu,
		cond: sync.NewCond(mu),
		kill: make(map[uint32]chan bool),
	}
}

func (que *taskQueue) add(items ...*pb.SimulationRun) {
	que.cond.L.Lock()
	defer que.cond.L.Unlock()

	que.tasks = append(que.tasks, items...)
	que.size = uint32(len(que.tasks))
	que.cond.Broadcast()
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

func (que *taskQueue) killLingering() {
	que.mu.Lock()
	defer que.mu.Unlock()

	for id, ch := range que.kill {
		ch <- true
		delete(que.kill, id)
	}
}

func (que *taskQueue) linger() (linger bool) {

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

	live := make(chan bool)
	defer close(live)

	go func() {
		que.cond.L.Lock()
		que.cond.Wait()
		live <- true
		que.cond.L.Unlock()
	}()

	select {
	case <-live:
		return true
	case <-die:
		return false
	}
}
