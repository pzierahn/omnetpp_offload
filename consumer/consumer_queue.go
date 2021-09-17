package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"sync"
	"sync/atomic"
)

type queue struct {
	cond  *sync.Cond
	size  int32
	tasks []*pb.SimulationRun
}

func newQueue() (que *queue) {
	return &queue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (que *queue) add(items ...*pb.SimulationRun) {
	que.cond.L.Lock()
	defer que.cond.L.Unlock()

	que.tasks = append(que.tasks, items...)
	atomic.SwapInt32(&que.size, int32(len(que.tasks)))
	que.cond.Broadcast()
}

func (que *queue) pop() (item *pb.SimulationRun, ok bool) {
	que.cond.L.Lock()
	defer que.cond.L.Unlock()

	if len(que.tasks) == 0 {
		return
	}

	ok = true
	item, que.tasks = que.tasks[0], que.tasks[1:]
	atomic.SwapInt32(&que.size, int32(len(que.tasks)))
	que.cond.Broadcast()

	return
}

func (que *queue) len() (size int32) {
	size = atomic.LoadInt32(&que.size)
	return
}

func (que *queue) onUpdate(callback func() (cancel bool)) {

	cancel := callback()

	for !cancel {
		que.cond.L.Lock()
		que.cond.Wait()

		cancel = callback()

		que.cond.L.Unlock()
	}

	return
}
