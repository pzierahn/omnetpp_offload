package consumer

import (
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

type consumer struct {
	config     *Config
	simulation *pb.Simulation
	connMu     sync.RWMutex

	finished sync.WaitGroup
	//allocCond *sync.Cond
	//allocate  []*pb.SimulationRun
	allocate *queue
	//allocator chan *pb.SimulationRun

	// TODO: Persist bytes to HD
	simulationSource []byte
	binaries         map[string][]byte
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Consumer ")
}
