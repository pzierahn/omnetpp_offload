package distribute

import (
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"path/filepath"
)

type Config struct {
	omnetpp.Config
	Tag             string   `json:"tag"`
	SimulateConfigs []string `json:"run"`
	Exclude         []string `json:"exclude"`
	SimulationId    string   `json:"-"`
}

func (config *Config) generateId() {

	tag := config.Tag

	if tag == "" {
		tag = filepath.Base(config.Path)
	}

	config.SimulationId = simple.NamedId(tag, 8)

	return
}
