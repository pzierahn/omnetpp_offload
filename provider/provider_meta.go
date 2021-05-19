package provider

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/grpc/metadata"
	"runtime"
)

type Meta struct {
	ProviderId string `json:"providerId"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	NumCPUs    int    `json:"numCPUs"`
}

func (info Meta) MarshallMeta() (md metadata.MD) {

	md = metadata.New(map[string]string{
		"providerId": info.ProviderId,
		"os":         info.Os,
		"arch":       info.Arch,
		"numCPUs":    fmt.Sprint(info.NumCPUs),
	})

	return
}

func (info *Meta) UnMarshallMeta(md metadata.MD) {

	info.ProviderId = utils.MetaStringFallback(md, "providerId", "")
	info.Os = utils.MetaStringFallback(md, "os", "")
	info.Arch = utils.MetaStringFallback(md, "arch", "")
	info.NumCPUs = utils.MetaIntFallback(md, "numCPUs", 0)

	return
}

func NewDeviceInfo(workerId string) (info Meta) {

	info = Meta{
		ProviderId: workerId,
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		NumCPUs:    runtime.NumCPU(),
	}

	return
}
