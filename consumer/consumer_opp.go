package consumer

import (
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
)

func extractConfigs(omnet omnetpp.OmnetProject) (configs map[string][]string, err error) {

	oppConfigs, err := omnet.GetConfigs()
	if err != nil {
		err = fmt.Errorf("couldn't get simulation configs: %v", err)
		return
	}

	configs = make(map[string][]string)

	for _, configName := range oppConfigs {

		var numbers []string
		numbers, err = omnet.GetRunNumbers(configName)
		if err != nil {
			return
		}

		configs[configName] = numbers
	}

	return
}
