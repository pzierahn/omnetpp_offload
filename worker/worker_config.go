package worker

type Config struct {
	WorkerId      string `json:"workerId,omitempty"`
	DeviceName    string `json:"deviceName,omitempty"`
	BrokerAddress string `json:"brokerAddress,omitempty"`
}
