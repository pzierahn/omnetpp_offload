package simulation

import (
	"github.com/patrickz98/project.go.omnetpp/simple"
)

type Config struct {
	Id      string
	Name    string
	Path    string
	Configs []string
}

func New(sourcePath, name string) (config Config) {
	config.Id = simple.NamedId(name, 6)
	config.Name = name
	config.Path = sourcePath

	return
}
