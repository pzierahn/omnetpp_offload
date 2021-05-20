package sysinfo

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetUtilization() (utilization *pb.Utilization, err error) {

	memo, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	utilization = &pb.Utilization{
		CpuUsage:    float32(GetCPUUsage()),
		MemoryTotal: memo.Total,
		MemoryUsed:  memo.Used,
		Updated:     timestamppb.Now(),
	}

	return
}
