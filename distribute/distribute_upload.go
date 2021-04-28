package distribute

import (
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

func Upload(conn gconfig.GRPCConnection, config *Config) (ref *pb.StorageRef, err error) {

	logger.Println("zipping", config.Path)

	buf, err := simple.TarGz(config.Path, config.SimulationId, config.Exclude...)
	if err != nil {
		return
	}

	logger.Printf("uploading %s to %s\n", config.SimulationId, conn.DialAddr())

	store := storage.InitClient(conn)
	ref, err = store.Upload(&buf, storage.FileMeta{
		Bucket:   config.SimulationId,
		Filename: "source.tar.gz",
	})

	return
}
