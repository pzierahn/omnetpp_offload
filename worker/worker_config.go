package worker

import "github.com/patrickz98/project.go.omnetpp/common"

type Config struct {
	workerId   string
	WorkerName string                `json:"workerName,omitempty"`
	Broker     common.GRPCConnection `json:"broker,omitempty"`
	DevoteCPUs int                   `json:"devoteCPUs,omitempty"`
}
