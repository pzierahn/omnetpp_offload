package storage

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
)

type Server struct {
	pb.UnimplementedStorageServer
}
