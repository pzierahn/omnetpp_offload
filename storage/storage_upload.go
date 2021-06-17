package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc/metadata"
	"io"
	"time"
)

func (client *Client) Upload(data io.Reader, meta FileMeta) (ref *pb.StorageRef, err error) {

	md := metadata.New(map[string]string{
		"bucket":   meta.Bucket,
		"filename": meta.Filename,
	})

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.storage.Push(ctx)
	if err != nil {
		return
	}

	start := time.Now()

	for chunk := range streamReader(data) {
		parcel := pb.StorageParcel{
			Size:    int32(chunk.size),
			Offset:  int64(chunk.offset),
			Payload: chunk.payload,
		}

		err = stream.Send(&parcel)
		if err != nil {
			return
		}

		if chunk.offset%8 == 0 {
			log.Printf("uploaded %s (%v)",
				simple.ByteSize(uint64(chunk.size+chunk.offset)), meta.Filename)
		}
	}

	ref, err = stream.CloseAndRecv()
	if err != nil {
		return
	}

	log.Printf("upload %v in %v\n", meta, time.Now().Sub(start))

	return
}
