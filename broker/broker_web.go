package broker

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/stargate"
)

func (broker *broker) providerUtilization() (overview *pb.Utilizations) {
	broker.pmu.RLock()
	defer broker.pmu.RUnlock()

	overview = &pb.Utilizations{
		Providers: make(map[string]*pb.Utilization),
	}

	broker.umu.RLock()
	for id, utilization := range broker.utilization {
		overview.Providers[id] = utilization
	}
	broker.umu.RUnlock()

	return
}

func (broker *broker) startDebugWebAPI() {

	simple.Watch("/providers", func() interface{} {
		return broker.providerUtilization()
	})

	simple.Watch("/stargate", func() interface{} {
		return stargate.DebugValues()
	})

	simple.StartWatchServer("")
}
