package storage

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net"
	"os"
)

const (
	port     = ":50052"
	dataPath = "data/storage"
)

type storage struct {
	pb.UnimplementedStorageServer
}

func (server *storage) Get(req *pb.StorageRef, stream pb.Storage_GetServer) (err error) {

	bucket := req.GetBucket()
	filename := req.GetFilename()
	filepath := dataPath + "/" + bucket + "/" + filename

	logger.Println("get", bucket, filename)

	file := fileStream(filepath)

	packages := 0

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

		packages++
	}

	logger.Println("packages send", packages)

	return
}

func (server *storage) Put(stream pb.Storage_PutServer) (err error) {

	var filepath string
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

	filepath = bucket + "/" + filename
	dataFile := dataPath + "/" + filepath

	_ = os.MkdirAll(dataPath+"/"+bucket, 0755)

	logger.Println("new upload request", bucket, "-->", filename)

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterStorageServer(server, &storage{})
	if err = server.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
