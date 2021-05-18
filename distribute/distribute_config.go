package distribute

import (
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
)

type Config struct {
	omnetpp.Config
	Tag             string   `json:"tag"`
	SimulateConfigs []string `json:"run"`
	Exclude         []string `json:"exclude"`
}
