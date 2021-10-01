package eval

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"time"
)

const (
	ConnectLocal = "Local"
	ConnectP2P   = "P2P"
	ConnectRelay = "Relay"
)

func LogSetup(connect string, details *pb.ProviderInfo) {

	timestamp := time.Now()
	ts, _ := timestamp.MarshalText()

	_, _ = cli.Setup(context.Background(), &pb.SetupEvent{
		TimeStamp:  string(ts),
		ProviderId: details.ProviderId,
		Connect:    connect,
		NumCPUs:    details.NumCPUs,
		Arch:       details.Arch.Arch,
		Os:         details.Arch.Os,
	})
}
