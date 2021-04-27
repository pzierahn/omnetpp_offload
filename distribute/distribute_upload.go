package distribute

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

func Upload(config *Config) (ref *pb.StorageRef, err error) {

	logger.Println("zipping", config.Path)

	buf, err := simple.TarGz(config.Path, config.SimulationId, config.Exclude...)
	if err != nil {
		return
	}

	logger.Println("uploading", config.SimulationId)

	ref, err = storage.Upload(&buf, storage.FileMeta{
		Bucket:   config.SimulationId,
		Filename: "source.tar.gz",
	})

	return
}
