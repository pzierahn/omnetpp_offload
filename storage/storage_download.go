package storage

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"io"
	"time"
)

func Download(file *pb.StorageRef) (byt io.Reader, err error) {

	conn, err := grpc.Dial(storageAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewStorageClient(conn)
	stream, err := client.Get(context.Background(), file)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	start := time.Now()

	packages := 0

	for {
		var parcel *pb.StorageParcel
		parcel, err = stream.Recv()

		if err == io.EOF {
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

	byt = &buf

	return
}
