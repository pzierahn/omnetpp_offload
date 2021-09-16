package storage

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"io"
	"time"
)

func (client *Client) Download(ctx context.Context, file *pb.StorageRef) (byt []byte, err error) {

	stream, err := client.storage.Pull(ctx, file)
	if err != nil {
		return
	}

	start := time.Now()
	packages := 0

	var buf bytes.Buffer

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

	log.Printf("Download: %s size=%d packages=%d time=%v",
		file.Filename, buf.Len(), packages, time.Now().Sub(start))

	byt = buf.Bytes()

	return
}
