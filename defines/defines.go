package defines

import (
	"os"
	"path/filepath"
)

const (
	DefaultPort = 50051
	DataPath    = "data"
	configDir   = "omnetpp-edge"
)

var (
	SimulationPath = filepath.Join(DataPath, "simulations")
)

func ConfigDir() (dir string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	dir = filepath.Join(dir, configDir)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	return
}
