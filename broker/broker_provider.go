package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	prov "github.com/patrickz98/project.go.omnetpp/provider"
	"sync"
)

func assignId(assignment *pb.Assignment) (id string) {

	switch work := assignment.Do.(type) {
	case *pb.Assignment_Build:
		id = fmt.Sprintf("%s.compile", work.Build.SimulationId)

	case *pb.Assignment_Run:
		id = fmt.Sprintf("%s.%s.%s",
			work.Run.SimulationId, work.Run.Config, work.Run.RunNumber)
	}

	return
}

type provider struct {
	sync.RWMutex
	id          string
	arch        *pb.Arch
	numCPUs     uint32
	utilization *pb.Utilization
	building    string
	assignments map[taskId]*pb.SimulationRun
	assign      chan *pb.Assignment
	//listener    map[chan<- *pb.ProviderState]interface{}
}

func newProvider(meta prov.Meta) (node *provider, err error) {

	id := meta.ProviderId
	if id == "" {
		err = fmt.Errorf("missing providerId in metadata")
		return
	}

	os := meta.Os
	if os == "" {
		err = fmt.Errorf("missing os in metadata")
		return
	}

	arch := meta.Arch
	if arch == "" {
		err = fmt.Errorf("missing arch in metadata")
		return
	}

	numCPUs := meta.NumCPUs
	if numCPUs == 0 {
		err = fmt.Errorf("missing numCPUs in metadata")
		return
	}

	node = &provider{
		id: id,
		arch: &pb.Arch{
			Os:   os,
			Arch: arch,
		},
		numCPUs:     uint32(numCPUs),
		assignments: make(map[taskId]*pb.SimulationRun),
		assign:      make(chan *pb.Assignment),
	}

	return
}

func (node *provider) assignRun(assignment *pb.SimulationRun) {
	node.Lock()
	defer node.Unlock()

	id := tId(assignment)
	node.assignments[id] = assignment

	node.assign <- &pb.Assignment{Do: &pb.Assignment_Run{
		Run: assignment,
	}}
}

func (node *provider) assignCompile(assignment *pb.Build) {
	node.Lock()
	defer node.Unlock()

	node.building = assignment.SimulationId

	node.assign <- &pb.Assignment{Do: &pb.Assignment_Build{
		Build: assignment,
	}}
}

func (node *provider) setUtilization(utilization *pb.Utilization) {
	node.Lock()
	defer node.Unlock()

	node.utilization = utilization
}

func (node *provider) busy() (busy bool) {
	node.RLock()
	defer node.RUnlock()

	if node.building != "" {
		busy = true
		return
	}

	busy = (node.utilization == nil) ||
		(node.utilization.CpuUsage >= 50.0) ||
		(len(node.assignments) >= int(node.numCPUs))

	return
}

func (node *provider) freeSlots() (num int) {
	node.RLock()
	defer node.RUnlock()

	num = int(node.numCPUs) - len(node.assignments)

	return
}

func (node *provider) close() {
	node.Lock()
	defer node.Unlock()

	close(node.assign)
}

type providerManager struct {
	sync.RWMutex
	provider map[string]*provider
}

func newProviderManager() (pm providerManager) {
	return providerManager{
		provider: make(map[string]*provider),
	}
}

func (pm *providerManager) add(node *provider) {

	logger.Printf("providerManager: add id=%v", node.id)

	pm.Lock()
	defer pm.Unlock()

	pm.provider[node.id] = node
}

func (pm *providerManager) remove(node *provider) {

	logger.Printf("providerManager: remove id=%v", node.id)

	pm.Lock()
	defer pm.Unlock()

	delete(pm.provider, node.id)
}
