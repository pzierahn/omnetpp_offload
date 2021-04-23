package simulation

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

func Upload(config Config) (ref *pb.StorageRef, err error) {

	logger.Println("zipping", config.Path)

	buf, err := simple.TarGz(config.Path, config.Id)
	if err != nil {
		return
	}

	logger.Println("uploading", config.Id)

	ref, err = storage.Upload(&buf, storage.FileMeta{
		Bucket:   config.Id,
		Filename: "source.tar.gz",
	})
	if err != nil {
		return
	}

	return
}
