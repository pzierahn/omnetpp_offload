package omnetpp

import (
	"os/exec"
	"regexp"
	"strings"
)

func (project *OmnetProject) GetConfigs() (configs []string, err error) {

	//
	// Get configs
	//

	simConfigs := exec.Command("./"+project.simulationExe, "-s", "-a")
	simConfigs.Dir = project.SourcePath

	byt, err := simConfigs.CombinedOutput()
	if err != nil {
		return
	}

	output := string(byt)
	output = strings.TrimSpace(output)

	reg := regexp.MustCompile(`Config (.+?):`)
	matches := reg.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		configs = append(configs, match[1])
	}

	return
}

func (project *OmnetProject) GetRunNumbers(config string) (numbers []string, err error) {

	//
	// Get runnumbers
	//

	runnumbers := exec.Command("./"+project.simulationExe, "-c", config, "-s", "-q", "runnumbers")
	runnumbers.Dir = project.SourcePath

	byt, err := runnumbers.CombinedOutput()
	if err != nil {
		return
	}

	output := string(byt)
	output = strings.TrimSpace(output)
	numbers = strings.Split(output, " ")

	return
}
