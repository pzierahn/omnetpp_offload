package defines

import "path/filepath"

const (
	Port     = ":50051"
	Address  = "192.168.0.11" + Port
	DataPath = "data"
)

var (
	Simulation = filepath.Join(DataPath, "simulations")
)
