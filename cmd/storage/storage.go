package main

import (
	"flag"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

var server bool
var upload string
var download string

func init() {
	flag.BoolVar(&server, "server", false, "start storage server")
	flag.StringVar(&upload, "upload", "", "upload path")
	flag.StringVar(&download, "download", "", "download file")
}

func main() {

	flag.Parse()

	if server {
		storage.StartServer()
	}

	if upload != "" {
		//storage.Upload(upload)
	}
}
