package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
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
			logger.Fatalln(err)
		}
	}

	ref, err = stream.CloseAndRecv()
	if err != nil {
		return
	}

	logger.Printf("upload %v in %v\n", meta, time.Now().Sub(start))

	return
}
