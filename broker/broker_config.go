package broker

type Config struct {
	Port         int  `json:"port"`
	WebInterface bool `json:"webInterface"`
}
