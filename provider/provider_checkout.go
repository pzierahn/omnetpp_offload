package provider

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"os"
	"path/filepath"
)

func (client *workerConnection) checkout(source *pb.Source) (path string, err error) {

	buf, err := client.storage.Download(source.Source)
	if err != nil {
		return
	}

	path = filepath.Join(cachePath, source.SimulationId)

	err = simple.UnTarGz(cachePath, &buf)
	if err != nil {
		_ = os.RemoveAll(path)
	}

	return
}
