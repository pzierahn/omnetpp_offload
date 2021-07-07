package storage

import (
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"path/filepath"
)

func (server *Server) Push(stream pb.Storage_PushServer) (err error) {

	var filename string
	var bucket string

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		err = fmt.Errorf("metadata missing")
		log.Println(err)
		return
	}

	filename, err = utils.MetaString(md, "filename")
	if err != nil {
		log.Println(err)
		return
	}

	bucket, err = utils.MetaString(md, "bucket")
	if err != nil {
		log.Println(err)
		return
	}

	dataFile := filepath.Join(storagePath, bucket, filename)
	dir, _ := filepath.Split(dataFile)
	_ = os.MkdirAll(dir, 0755)

	log.Println("push", dataFile)

	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	for {
		var parcel *pb.StorageParcel
		parcel, err = stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		_, err = file.WriteAt(parcel.Payload, int64(parcel.Offset))
		if err != nil {
			return
		}
	}

	err = stream.SendAndClose(&pb.StorageRef{
		Bucket:   bucket,
		Filename: filename,
	})

	return
}
