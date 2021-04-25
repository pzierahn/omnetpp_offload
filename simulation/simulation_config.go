package simulation

import (
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"path/filepath"
)

type Config struct {
	omnetpp.Config
	Tag             string   `json:"tag"`
	SimulationId    string   `json:"-"`
	SimulateConfigs []string `json:"simulateConfigs"`
}

func (config *Config) GenerateId() {

	tag := config.Tag

	if tag == "" {
		tag = filepath.Base(config.Path)
	}

	config.SimulationId = simple.NamedId(tag, 8)

	return
}
