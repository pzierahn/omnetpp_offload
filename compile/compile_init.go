package compile

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

type Compiler struct {
	Broker         pb.BrokerClient
	Storage        storage.Client
	SimulationId   string
	SimulationBase string
	OppConfig      *pb.OppConfig
	cleanedFiles   map[string]bool
	compiledFiles  map[string]bool
}
