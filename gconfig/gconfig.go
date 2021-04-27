package gconfig

type Config struct {
	Broker GRPCConnection `json:"broker,omitempty"`
	Worker Worker         `json:"worker,omitempty"`
}

const (
	ParseBroker = iota
	ParseWorker
)

var ParseAll = []int{
	ParseBroker,
	ParseWorker,
}
