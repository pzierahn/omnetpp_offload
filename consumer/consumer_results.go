package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
)

// Extract Results files to the right place
func (cons *consumer) extractResults(byt []byte) {
	err := simple.ExtractTarGz(cons.config.Path, byt)
	if err != nil {
		log.Printf("ExtractTarGz failed: %v", err)
	}
}
