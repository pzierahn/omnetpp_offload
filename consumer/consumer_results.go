package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
)

// Extract Results files to the right place
func (sim *simulation) extractResults(byt []byte) {
	err := simple.ExtractTarGz(sim.config.Path, byt)
	if err != nil {
		log.Printf("cloudn't extract files: %v", err)
	}
}
