package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/grpc"
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
	ctx        context.Context
	config     *Config
	simulation *pb.Simulation
	connMu     sync.RWMutex
	bconn      *grpc.ClientConn

	finished sync.WaitGroup
	allocate *queue

	// TODO: Persist bytes to HD
	simulationSource []byte
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Consumer ")
}
