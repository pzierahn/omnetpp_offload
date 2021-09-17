package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
)

type Server struct {
	pb.UnimplementedStorageServer
}

// Delete is the gRPC server implementation for deleting files.
func (server *Server) Delete(_ context.Context, ref *pb.StorageRef) (res *emptypb.Empty, err error) {

	res = &emptypb.Empty{}

	log.Printf("Delete: %+v", ref)

	filename := filepath.Join(storagePath, ref.Bucket, ref.Filename)
	err = os.RemoveAll(filename)

	return
}

// Drop is the gRPC server implementation for deleting buckets.
func (server *Server) Drop(_ context.Context, ref *pb.BucketRef) (res *emptypb.Empty, err error) {

	res = &emptypb.Empty{}

	log.Printf("Drop: %+v", ref)

	filename := filepath.Join(storagePath, ref.Bucket)
	err = os.RemoveAll(filename)

	return
}
