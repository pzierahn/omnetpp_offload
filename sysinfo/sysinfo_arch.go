package sysinfo

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"runtime"
)

func Arch() (arch *pb.Arch) {
	arch = &pb.Arch{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	return
}

func ArchSignature() (sig string) {
	return Signature(Arch())
}

func Signature(arch *pb.Arch) (sig string) {
	sig = arch.Os + "_" + arch.Arch
	return
}
