package storage

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"io"
	"time"
)

func (client *Client) Download(ctx context.Context, file *pb.StorageRef) (buf bytes.Buffer, err error) {

	stream, err := client.storage.Pull(ctx, file)
	if err != nil {
		return
	}

	start := time.Now()
	packages := 0

	for {
		var parcel *pb.StorageParcel
		parcel, err = stream.Recv()

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return
		}

		_, err = buf.Write(parcel.Payload)
		if err != nil {
			return
		}

		packages++
	}

	log.Printf("received %d packages in %v", packages, time.Now().Sub(start))

	return
}
