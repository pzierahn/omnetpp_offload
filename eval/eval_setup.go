package eval

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
)

const (
	ConnectLocal = "Local"
	ConnectP2P   = "P2P"
	ConnectRelay = "Relay"
)

func LogSetup(connect string, details *pb.ProviderInfo) {

	//
	// TODO
	//

	//timestamp := time.Now()
	//ts, _ := timestamp.MarshalText()
	//
	//_, _ = cli.Setup(context.Background(), &pb.SetupEvent{
	//	TimeStamp:  string(ts),
	//	ProviderId: details.ProviderId,
	//	Connect:    connect,
	//	NumCPUs:    details.NumCPUs,
	//	NumJobs:    details.NumJobs,
	//	Arch:       details.Arch.Arch,
	//	Os:         details.Arch.Os,
	//})
}
