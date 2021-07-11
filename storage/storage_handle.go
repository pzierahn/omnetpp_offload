package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
)

func (server *Server) Delete(_ context.Context, ref *pb.StorageRef) (res *emptypb.Empty, err error) {

	res = &emptypb.Empty{}

	log.Printf("Delete: %+v", ref)

	filename := filepath.Join(storagePath, ref.Bucket, ref.Filename)
	err = os.RemoveAll(filename)

	return
}
