package storage

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (server *Server) Put(rea io.Reader, ref *pb.StorageRef) (err error) {

	log.Printf("put: %v", ref)

	dest := filepath.Join(storagePath, ref.Bucket, ref.Filename)

	base, _ := filepath.Split(dest)
	_ = os.MkdirAll(base, 0755)

	file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, rea)

	return
}

func (server *Server) Get(ref *pb.StorageRef) (byt []byte, err error) {

	log.Printf("get: %v", ref)

	src := filepath.Join(storagePath, ref.Bucket, ref.Filename)
	byt, err = ioutil.ReadFile(src)

	return
}
