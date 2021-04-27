package storage

import (
	"bytes"
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"io"
	"time"
)

func (client *Client) Download(file *pb.StorageRef) (buf bytes.Buffer, err error) {

	stream, err := client.storage.Pull(context.Background(), file)
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
			logger.Fatalln(err)
		}

		_, err = buf.Write(parcel.Payload)
		if err != nil {
			logger.Fatalln(err)
		}

		packages++
	}

	logger.Printf("received %d packages in %v\n", packages, time.Now().Sub(start))

	return
}
