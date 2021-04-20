package simulation

import (
	"com.github.patrickz98.omnet/simple"
	"time"
)

type Config struct {
	Id      string
	Name    string
	Created time.Time
	Path    string
}

func New(sourcePath, name string) (config Config) {
	config.Id = simple.NamedId(name, 6)
	config.Name = name
	config.Path = sourcePath
	config.Created = time.Now()

	return
}
