package main

import (
	"context"
	"flag"
	"github.com/patrickz98/project.go.omnetpp/defines"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/worker"
	"runtime"
)

const (
	configFile = "worker-config.json"
)

var (
	config     worker.Config
	workerName = simple.GetHostnameShort()
)

func init() {

	//if _, err := os.Stat(configPath); err == nil {
	//	_ = simple.UnmarshallFile(configPath, &config)
	//}

	//flag.BoolVar(&saveConfig, "saveConfig", false, "create a new config file")
	//flag.BoolVar(&loadConfig, "loadConfig", false, "set worker name")
	//flag.BoolVar(&loadConfig, "useDefaultConfig", false, "set worker name")

	flag.StringVar(&config.WorkerName, "workerName", workerName, "set worker name")
	flag.StringVar(&config.BrokerAddress, "brokerAddress", "", "set broker address")
	flag.IntVar(&config.BrokerPort, "brokerPort", defines.DefaultPort, "set broker port")
	flag.IntVar(&config.DevoteCPUs, "devoteCPUs", runtime.NumCPU(), "set number of CPU cores")
}

func main() {

	flag.Parse()

	//if saveConfig {
	//	byt, err := json.MarshalIndent(config, "", "  ")
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	file := filepath.Join(defines.ConfigDir(), configFile)
	//	fmt.Printf("saving config to %s\n", file)
	//
	//	err = ioutil.WriteFile(file, byt, 0644)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	return
	//}

	conn, err := worker.Connect(config)
	if err != nil {
		panic(err)
	}

	if err = conn.StartLink(context.Background()); err != nil {
		panic(err)
	}
}
