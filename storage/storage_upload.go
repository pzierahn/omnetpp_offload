package storage

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc/metadata"
	"sync/atomic"
	"time"
)

type UploadProgress struct {
	Parcels  uint32
	Uploaded uint64
	Percent  float32
}

// Upload uploads a file to the storage server and returns a storage reference.
func (client *Client) Upload(meta *File, fb chan<- UploadProgress) (ref *pb.StorageRef, err error) {

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

	if fb != nil {

		//
		// feedback upload progress
		//

		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		go func() {
			for range tick.C {
				upl := atomic.LoadUint64(&uploaded)

				fb <- UploadProgress{
					Parcels:  atomic.LoadUint32(&parcels),
					Uploaded: upl,
					Percent:  100 * (float32(upl) / float32(len(meta.Data))),
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
