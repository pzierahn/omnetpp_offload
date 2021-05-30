package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	"io/ioutil"
	"log"
	"path/filepath"
)

var path string
var configPath string

var clean bool
var makemake bool
var compile bool
var printConfigs bool

var configname string
var printRunNumbers bool
var runfilter string

var tgz string

func init() {

	//
	// Setup
	//

	flag.StringVar(&path, "path", ".", "simulation path")
	flag.StringVar(&configPath, "configfile", "opp-edge-config.json", "simulation config JSON")

	//
	// Actions
	//

	flag.BoolVar(&clean, "clean", false, "clean simulation binaries")
	flag.BoolVar(&makemake, "makemake", false, "create Makefile")
	flag.BoolVar(&compile, "compile", false, "compile the simulation")
	flag.BoolVar(&printConfigs, "printConfigs", false, "print simulation configurations")

	flag.StringVar(&configname, "configname", "General", "Select a configuration for execution")
	flag.BoolVar(&printRunNumbers, "printRunNumbers", false, "print run numbers for simulation configuration")
	flag.StringVar(&runfilter, "runfilter", "", "run configuration")
	flag.StringVar(&tgz, "tgz", "", "tgz results filename")
}

func main() {
	flag.Parse()

	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	var config omnetpp.Config
	config.Path = path

	byt, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(byt, &config)
	if err != nil {
		log.Panicln(err)
	}

	opp := omnetpp.New(&config)

	if clean {
		err = opp.Clean()
		if err != nil {
			panic(err)
		}
	}

	if makemake {
		err = opp.MakeMake()
		if err != nil {
			log.Panicln(err)
		}
	}

	if compile {
		err = opp.Compile()
		if err != nil {
			panic(err)
		}
	}

	if printConfigs {
		configs, err := opp.GetConfigs()
		if err != nil {
			panic(err)
		}

		fmt.Println(configs)
	}

	if printRunNumbers {
		numbers, err := opp.GetRunNumbers(configname)
		if err != nil {
			panic(err)
		}

		fmt.Println(numbers)
	}

	if runfilter != "" {
		err = opp.Run(configname, runfilter)
		if err != nil {
			panic(err)
		}
	}

	if tgz != "" {
		byt, err := opp.ZipResults()
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(tgz, byt.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}
