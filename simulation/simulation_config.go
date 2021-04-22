package simulation

import (
	"com.github.patrickz98.omnet/simple"
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
