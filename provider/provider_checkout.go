package provider

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"os"
	"path/filepath"
)

func (client *workerConnection) checkout(simulationId string) (path string, err error) {

	req := &pb.SimulationId{Id: simulationId}
	src, err := client.broker.GetSource(context.Background(), req)
	if err != nil {
		return
	}

	buf, err := client.storage.Download(src.Source)
	if err != nil {
		return
	}

	path = filepath.Join(cachePath, simulationId)

	err = simple.UnTarGz(cachePath, &buf)
	if err != nil {
		_ = os.RemoveAll(path)
	}

	return
}
