package simulation

import (
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
	"fmt"
)

func extractConfigs(omnet omnetpp.OmnetProject) (configs []*pb.Config, err error) {

	oppConfigs, err := omnet.GetConfigs()
	if err != nil {
		err = fmt.Errorf("couldn't get simulation configs: %v", err)
		return
	}

	for _, configName := range oppConfigs {

		var numbers []string
		numbers, err = omnet.GetRunNumbers(configName)
		if err != nil {
			return
		}

		pbConf := pb.Config{
			Name:       configName,
			RunNumbers: numbers,
		}

		configs = append(configs, &pbConf)
	}

	return
}
