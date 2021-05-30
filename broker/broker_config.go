package broker

type Config struct {
	StunPort     int  `json:"stunPort"`
	BrokerPort   int  `json:"brokerPort"`
	WebInterface bool `json:"webInterface"`
}
