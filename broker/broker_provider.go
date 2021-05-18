package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"sync"
)

type providerId string

func pId(pState *pb.ProviderState) (id providerId) {
	return providerId(pState.ProviderId)
}

type providerManager struct {
	sync.RWMutex
	provider map[providerId]*pb.ProviderState
	listener map[providerId]map[chan<- *pb.ProviderState]interface{}
	work     map[providerId]chan<- *pb.Work
}

func newProviderManager() (union providerManager) {
	return providerManager{
		provider: make(map[providerId]*pb.ProviderState),
		listener: make(map[providerId]map[chan<- *pb.ProviderState]interface{}),
		work:     make(map[providerId]chan<- *pb.Work),
	}
}

func (pm *providerManager) update(state *pb.ProviderState) {

	id := pId(state)

	// logger.Printf("update: %v (%.0f%%)", id, state.CpuUsage)

	pm.Lock()
	pm.provider[id] = state
	pm.Unlock()

	pm.RLock()
	//logger.Printf("update: send events to %v listeners", len(pm.listener[id]))

	for listener := range pm.listener[id] {
		listener <- state
	}
	pm.RUnlock()
}

func (pm *providerManager) remove(id providerId) {

	logger.Printf("remove: id=%v", id)

	pm.Lock()
	delete(pm.provider, id)
	delete(pm.listener, id)
	delete(pm.work, id)
	pm.Unlock()
}

func (pm *providerManager) newProvider(id providerId, ch chan *pb.Work) {
	pm.Lock()
	defer pm.Unlock()

	pm.work[id] = ch
}

func (pm *providerManager) removeWorker(id providerId) {
	pm.Lock()
	defer pm.Unlock()

	delete(pm.work, id)
}

func (pm *providerManager) newListener(id providerId, ch chan *pb.ProviderState) {
	pm.Lock()
	defer pm.Unlock()

	if pm.listener[id] == nil {
		pm.listener[id] = make(map[chan<- *pb.ProviderState]interface{})
	}

	pm.listener[id][ch] = ch
}

func (pm *providerManager) removeListener(id providerId, ch chan *pb.ProviderState) {
	pm.Lock()
	defer pm.Unlock()

	delete(pm.listener[id], ch)
}
