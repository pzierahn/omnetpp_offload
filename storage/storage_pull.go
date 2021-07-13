package storage

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"os"
	"path/filepath"
	"time"
)

func (server *Server) Pull(req *pb.StorageRef, stream pb.Storage_PullServer) (err error) {

	filename := filepath.Join(storagePath, req.Bucket, req.Filename)

	log.Println("pull", req.Bucket, req.Filename)

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
	start := time.Now()

	reader := streamReader(file)

	for chunk := range reader {
		parcel := pb.StorageParcel{
			Size:    uint32(chunk.size),
			Offset:  uint64(chunk.offset),
			Payload: chunk.payload,
		}

		err = stream.Send(&parcel)
		if err != nil {
			return
		}

		packages++

		log.Printf("package %s->%s send %0.2f%%",
			req.Bucket, req.Filename, 100.0*(float64(chunk.offset+chunk.size)/float64(stat.Size())))
	}

	log.Printf("%s->%s packges %d in %v", req.Bucket, req.Filename, packages, time.Now().Sub(start))

	return
}
