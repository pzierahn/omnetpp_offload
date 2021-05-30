package provider

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"runtime"
)

func (prov *provider) info() (info *pb.ProviderInfo) {

	info = &pb.ProviderInfo{
		ProviderId: prov.providerId,
		Os: &pb.OS{
			Os:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		NumCPUs: uint32(runtime.NumCPU()),
	}

	return
}
