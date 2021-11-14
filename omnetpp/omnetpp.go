package omnetpp

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
	"os"
)

type Config struct {
	*pb.OppConfig
	Path string `json:"-"`
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "OMNeT++ ", log.LstdFlags|log.Lshortfile)
}

type OmnetProject struct {
	*Config
}

func New(config *Config) (project OmnetProject) {

	//
	// TODO: Add default values to config
	//

	project = OmnetProject{
		config,
	}

	return
}
