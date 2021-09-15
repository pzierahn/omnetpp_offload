package storage

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/grpc/metadata"
	"sync/atomic"
	"time"
)

type UploadInfo struct {
	Parcels  uint32
	Uploaded uint64
}

func (client *Client) Upload(meta *FileMeta, ch chan<- UploadInfo) (ref *pb.StorageRef, err error) {

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

	var parcels uint32
	var uploaded uint64

	if ch != nil {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		go func() {
			for range tick.C {
				ch <- UploadInfo{
					Parcels:  atomic.LoadUint32(&parcels),
					Uploaded: atomic.LoadUint64(&uploaded),
				}
			}
		}()
	}

	for chunk := range streamReader(bytes.NewReader(meta.Data)) {
		parcel := pb.StorageParcel{
			Size:    uint32(chunk.size),
			Offset:  uint64(chunk.offset),
			Payload: chunk.payload,
		}

		err = stream.Send(&parcel)
		if err != nil {
			return
		}

		atomic.AddUint32(&parcels, 1)
		atomic.AddUint64(&uploaded, uint64(chunk.size))
	}

	return stream.CloseAndRecv()
}
