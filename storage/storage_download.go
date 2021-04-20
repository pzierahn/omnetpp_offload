package storage

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"google.golang.org/grpc"
	"io"
	"os"
	"time"
)

func Download(filename string) {
	logger.Println("download", filename)

	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewStorageClient(conn)

	request := &pb.StorageRef{
		Bucket:   "tictic-1234",
		Filename: filename,
	}

	stream, err := client.Get(context.Background(), request)
	if err != nil {
		logger.Fatalln(err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Fatalln(err)
	}
	defer func() {
		_ = file.Close()
	}()

	start := time.Now()

	for {
		parcel, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Fatalln(err)
		}

		_, err = file.WriteAt(parcel.Payload, parcel.Offset)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	logger.Println("time", time.Now().Sub(start))
}
