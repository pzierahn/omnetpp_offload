package consumer

import (
	"bytes"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"sync"
)

type consumer struct {
	config      *Config
	simulation  *pb.Simulation
	connMu      sync.RWMutex
	connections map[string]*connection

	// TODO: Persist bytes to HD
	simulationTgz bytes.Buffer
	binaries      map[string]bytes.Buffer
}
