package gconfig

type Worker struct {
	Name       string `json:"name,omitempty"`
	DevoteCPUs int    `json:"devoteCPUs,omitempty"`
}
