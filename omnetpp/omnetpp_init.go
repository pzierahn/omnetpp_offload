package omnetpp

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
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
	project = OmnetProject{
		config,
	}

	return
}
