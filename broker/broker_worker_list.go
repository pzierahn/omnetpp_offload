package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"sync"
)

type workerList struct {
	sync.RWMutex
	workers map[string]*pb.ResourceCapacity
}

func (list *workerList) Get(key string) (info *pb.ResourceCapacity, ok bool) {
	list.RLock()

	info, ok = list.workers[key]

	list.RUnlock()

	return
}

func (list *workerList) Put(key string, value *pb.ResourceCapacity) {
	list.Lock()

	list.workers[key] = value

	list.Unlock()
}

func (list *workerList) Remove(key string) {
	list.Lock()

	delete(list.workers, key)

	list.Unlock()
}

func initWorkerList() (list workerList) {
	list = workerList{
		RWMutex: sync.RWMutex{},
		workers: make(map[string]*pb.ResourceCapacity),
	}

	return
}
