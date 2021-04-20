package storage

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func Upload(filename string) {
	logger.Println("upload", filename)

	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewStorageClient(conn)

	md := metadata.New(map[string]string{
		"bucket":   "tictic-1234",
		"filename": filename,
	})

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.Put(ctx)
	if err != nil {
		logger.Fatalf("connecting stream failed: %v", err)
	}

	start := time.Now()

	file := fileStream(filename)

	for chunk := range file {
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

	reply, err := stream.CloseAndRecv()
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("reply", reply.GetBucket(), reply.GetFilename())
	logger.Println("time", time.Now().Sub(start))
}
