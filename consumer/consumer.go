package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
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
	ctx      context.Context
	config   *Config
	id       string
	finished sync.WaitGroup
	queue    *taskQueue
	amu      sync.Mutex
	archLock map[string]*sync.Mutex
	bmu      sync.RWMutex
	binaries map[string][]byte

	// TODO: Persist bytes to HD
	source []byte
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
