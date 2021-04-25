package omnetpp

import (
	"bytes"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"path/filepath"
)

func (project *OmnetProject) ZipResults() (buf bytes.Buffer, err error) {

	resultsPath := filepath.Join(project.Path, project.ResultsPath)
	buf, err = simple.TarGz(resultsPath, "results")

	return
}
