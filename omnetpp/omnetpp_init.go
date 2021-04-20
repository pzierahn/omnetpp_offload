package omnetpp

import (
	"log"
	"os"
)

var logger *log.Logger

const (
	omnetBin = "/Users/patrick/Desktop/omnetpp-5.6.2/bin"
)

func init() {
	logger = log.New(os.Stderr, "Omnetpp", log.LstdFlags|log.Lshortfile)
}

type OmnetProject struct {
	SourcePath    string
	simulationExe string
}

func New(path string) (project OmnetProject) {
	project = OmnetProject{
		SourcePath:    path,
		simulationExe: "simulation",
	}

	return
}
