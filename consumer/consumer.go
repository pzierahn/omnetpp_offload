package consumer

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/omnetpp"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
	"sync"
)

type Config struct {
	omnetpp.Config
	Tag             string   `json:"tag"`
	SimulateConfigs []string `json:"run"`
	Ignore          []string `json:"ignore"`
}

type simulation struct {
	id       string
	ctx      context.Context
	config   *Config
	finished sync.WaitGroup
	queue    *taskQueue
	amu      sync.Mutex
	archLock map[string]*sync.Mutex
	bmu      sync.RWMutex
	binaries map[string][]byte
	onInit   chan uint32
	// TODO: Persist bytes to HD
	source []byte

	// TODO: Store connections central
	//cmu         sync.RWMutex
	//connections map[string]*providerConnection
}

func (sim *simulation) proto() (simulation *pb.Simulation) {
	return &pb.Simulation{
		Id:        sim.id,
		OppConfig: sim.config.OppConfig,
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Consumer ")
}
