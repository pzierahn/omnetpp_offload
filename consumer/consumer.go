package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"sync"
)

type Config struct {
	omnetpp.Config
	Tag             string   `json:"tag"`
	SimulateConfigs []string `json:"run"`
	Ignore          []string `json:"ignore"`
}

type consumer struct {
	consumerId  string
	config      *Config
	simulation  *pb.Simulation
	connMu      sync.RWMutex
	connections map[string]*connection

	finished  sync.WaitGroup
	allocCond *sync.Cond
	allocate  []*pb.SimulationRun
	allocator chan *pb.SimulationRun

	// TODO: Persist bytes to HD
	simulationTgz []byte
	binaries      map[string][]byte
}
