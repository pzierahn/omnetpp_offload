package provider

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/sysinfo"
	"runtime"
)

func (prov *provider) info() (info *pb.ProviderInfo) {

	info = &pb.ProviderInfo{
		ProviderId: prov.providerId,
		Arch:       sysinfo.Arch(),
		NumCPUs:    uint32(runtime.NumCPU()),
		NumJobs:    uint32(prov.numJobs),
	}

	return
}
