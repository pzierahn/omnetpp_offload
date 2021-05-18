package distribute

import (
	"bytes"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"runtime"
)

func UploadSource(server *grpc.ClientConn, simulationId string, config *Config) (ref *pb.StorageRef, err error) {

	logger.Println("zipping", config.Path)

	buf, err := simple.TarGz(config.Path, simulationId, config.Exclude...)
	if err != nil {
		return
	}

	logger.Printf("uploading source %s to %s\n", simulationId, server.Target())

	store := storage.ConnectClient(server)
	ref, err = store.Upload(&buf, storage.FileMeta{
		Bucket:   simulationId,
		Filename: "source.tar.gz",
	})

	return
}

func UploadBinary(server *grpc.ClientConn, simulationId string, buf *bytes.Buffer) (ref *pb.StorageRef, err error) {

	logger.Printf("uploading binary %s to %s\n", simulationId, server.Target())

	store := storage.ConnectClient(server)
	ref, err = store.Upload(buf, storage.FileMeta{
		Bucket:   simulationId,
		Filename: fmt.Sprintf("binary_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH),
	})

	return
}
