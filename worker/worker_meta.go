package worker

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/grpc/metadata"
	"runtime"
)

type DeviceInfo struct {
	WorkerId string `json:"workerId"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	NumCPUs  int    `json:"numCPUs"`
}

func (info DeviceInfo) MarshallMeta() (md metadata.MD) {

	md = metadata.New(map[string]string{
		"workerId": info.WorkerId,
		"os":       info.Os,
		"arch":     info.Arch,
		"numCPUs":  fmt.Sprint(info.NumCPUs),
	})

	return
}

func (info *DeviceInfo) UnMarshallMeta(md metadata.MD) {

	info.WorkerId = utils.MetaStringFallback(md, "workerId", "")
	info.Os = utils.MetaStringFallback(md, "os", "")
	info.Arch = utils.MetaStringFallback(md, "arch", "")
	info.NumCPUs = utils.MetaIntFallback(md, "numCPUs", 0)

	return
}

func NewDeviceInfo(workerId string) (info DeviceInfo) {

	info = DeviceInfo{
		WorkerId: workerId,
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCPUs:  runtime.NumCPU(),
	}

	return
}
