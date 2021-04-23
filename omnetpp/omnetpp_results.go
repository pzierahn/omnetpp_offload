package omnetpp

import (
	"bytes"
	"com.github.patrickz98.omnet/simple"
	"path/filepath"
)

func (project *OmnetProject) ZipResults() (buf bytes.Buffer, err error) {

	resultsPath := filepath.Join(project.SourcePath, "results")
	buf, err = simple.TarGz(resultsPath, "results")

	return
}
