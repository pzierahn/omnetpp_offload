package sysinfo

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"runtime"
)

func Arch() (arch *pb.Arch) {
	arch = &pb.Arch{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	return
}

func ArchSignature() (arch string) {
	arch = runtime.GOOS + "_" + runtime.GOARCH
	return
}
