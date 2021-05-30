package storage

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"os"
	"path/filepath"
)

func (server *Server) Pull(req *pb.StorageRef, stream pb.Storage_PullServer) (err error) {

	filename := filepath.Join(storagePath, req.Bucket, req.Filename)

	logger.Println("pull", req.Bucket, req.Filename)

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	var packages int

	reader := streamReader(file)

	for chunk := range reader {
		parcel := pb.StorageParcel{
			Size:    int32(chunk.size),
			Offset:  int64(chunk.offset),
			Payload: chunk.payload,
		}

		err = stream.Send(&parcel)
		if err != nil {
			logger.Fatalln(err)
		}

		packages++

		logger.Printf("packages %s->%s send %0.2f%%\n", req.Bucket, req.Filename, 100.0*(float64(chunk.offset+chunk.size)/float64(stat.Size())))
	}

	logger.Println("packages send", packages)

	return
}
