package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"io/ioutil"
)

var ref pb.StorageRef

var pull bool
var ls bool
var rm bool

var config gconfig.Config

func init() {
	flag.StringVar(&ref.Bucket, "bucket", "", "bucket name")
	flag.StringVar(&ref.Filename, "file", "", "filename")

	flag.BoolVar(&pull, "pull", false, "download file")
	flag.BoolVar(&ls, "ls", false, "show list files in bucket")
	flag.BoolVar(&rm, "rm", false, "delete file / bucket")

	config = gconfig.SourceAndParse(gconfig.ParseBroker)
}

func main() {

	flag.Parse()

	fmt.Printf("connecting to %s\n", config.Broker.DialAddr())

	client := storage.InitClient(config.Broker)
	defer client.Close()

	if ls {
		list, err := client.List(&ref)
		if err != nil {
			panic(err)
		}

		jsByt, _ := json.MarshalIndent(list, "", "  ")
		fmt.Println(string(jsByt))
	}

	if pull {
		byt, err := client.Download(&ref)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(ref.Filename, byt.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}

	if rm {
		res, err := client.Delete(&ref)
		if err != nil {
			panic(err)
		}

		jsByt, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(jsByt))
	}
}
