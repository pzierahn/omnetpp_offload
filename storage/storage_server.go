package storage

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net"
	"os"
	"path/filepath"
)

type storage struct {
	pb.UnimplementedStorageServer
}

func (server *storage) Get(req *pb.StorageRef, stream pb.Storage_GetServer) (err error) {

	filename := filepath.Join(storagePath, req.Bucket, req.Filename)

	logger.Println("get", req.Bucket, req.Filename)

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	var packages int

	reader := streamReader(file)

	for chunk := range reader {
		parcel := pb.StorageParcel{
			Size:    int32(chunk.size),
			Offset:  int64(chunk.offset),
			Payload: chunk.payload,
		}

		err = stream.Send(&parcel)
		if err != nil {
			logger.Fatalln(err)
		}

		packages++

		logger.Printf("packages %s->%s send %0.2f%%\n", req.Bucket, req.Filename, 100.0*(float64(chunk.offset+chunk.size)/float64(stat.Size())))
	}

	logger.Println("packages send", packages)

	return
}

func (server *storage) Put(stream pb.Storage_PutServer) (err error) {

	var filename string
	var bucket string

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		err = fmt.Errorf("metadata missing")
		logger.Println(err)
		return
	}

	filename, err = utils.MetaString(md, "filename")
	if err != nil {
		logger.Println(err)
		return
	}

	bucket, err = utils.MetaString(md, "bucket")
	if err != nil {
		logger.Println(err)
		return
	}

	dataFile := filepath.Join(storagePath, bucket, filename)

	_ = os.MkdirAll(filepath.Join(storagePath, bucket), 0755)

	logger.Println("put", bucket, "-->", filename)

	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	for {
		var parcel *pb.StorageParcel
		parcel, err = stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		_, err = file.WriteAt(parcel.Payload, parcel.Offset)
		if err != nil {
			return
		}
	}

	err = stream.SendAndClose(&pb.StorageRef{
		Bucket:   bucket,
		Filename: filename,
	})

	return
}

func StartServer() {

	logger.Println("start storage server")

	lis, err := net.Listen("tcp", storageAddress)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterStorageServer(server, &storage{})
	if err = server.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
