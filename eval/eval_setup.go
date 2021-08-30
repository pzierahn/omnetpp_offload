package eval

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"sync"
	"time"
)

const (
	ConnectLocal = "Local"
	ConnectP2P   = "P2P"
	ConnectRelay = "Relay"
)

type Setup struct {
	ScenarioId   string    `csv:"scenario_id"`
	SimulationId string    `csv:"simulation_id"`
	TrailId      string    `csv:"trail_id"`
	TimeStamp    time.Time `csv:"time_stamp"`
	Connect      string    `csv:"connect"`
	ProviderId   string    `csv:"provider_id"`
	NumCPUs      int       `csv:"num_cpus"`
	Arch         string    `csv:"arch"`
	Os           string    `csv:"os"`
}

var smu sync.Mutex
var setup []Setup

func LogSetup(connect string, details *pb.ProviderInfo) {
	smu.Lock()
	setup = append(setup, Setup{
		ScenarioId:   ScenarioId,
		SimulationId: SimulationId,
		TrailId:      TrailId,
		TimeStamp:    time.Now(),
		Connect:      connect,
		ProviderId:   details.ProviderId,
		NumCPUs:      int(details.NumCPUs),
		Arch:         details.Arch.Arch,
		Os:           details.Arch.Os,
	})
	smu.Unlock()
}
