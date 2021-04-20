package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"sync"
)

type workerList struct {
	sync.RWMutex
	workers map[string]*pb.ClientInfo
}

func (list *workerList) Get(key string) (info *pb.ClientInfo, ok bool) {
	list.RLock()

	info, ok = list.workers[key]

	list.RUnlock()

	return
}

func (list *workerList) Put(key string, value *pb.ClientInfo) {
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
		workers: make(map[string]*pb.ClientInfo),
	}

	return
}
