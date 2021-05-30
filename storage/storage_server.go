package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Server struct {
	pb.UnimplementedStorageServer
}

func (server *Server) List(_ context.Context, req *pb.StorageRef) (list *pb.StorageList, err error) {

	logger.Println("list", req.Bucket)

	dirname := filepath.Join(storagePath, req.Bucket)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return
	}

	list = &pb.StorageList{}
	for _, file := range files {
		list.Files = append(list.Files, file.Name())
	}

	return
}

func (server *Server) Delete(_ context.Context, req *pb.StorageRef) (res *pb.StorageStatus, err error) {

	logger.Println("delete", req.Bucket)

	dirname := filepath.Join(storagePath, req.Bucket, req.Filename)
	err = os.RemoveAll(dirname)

	res = &pb.StorageStatus{
		Deleted: err == nil,
	}

	return
}
