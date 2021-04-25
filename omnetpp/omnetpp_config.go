package omnetpp

import (
	"regexp"
	"strings"
)

func (project *OmnetProject) GetConfigs() (configs []string, err error) {

	//
	// Get configs
	//

	sim, err := project.command("-s", "-a")
	if err != nil {
		return
	}

	byt, err := sim.CombinedOutput()
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

	sim, err := project.command("-c", config, "-q", "runnumbers", "-s")
	if err != nil {
		return
	}

	byt, err := sim.CombinedOutput()
	if err != nil {
		return
	}

	output := string(byt)
	output = strings.TrimSpace(output)
	numbers = strings.Split(output, " ")

	return
}
