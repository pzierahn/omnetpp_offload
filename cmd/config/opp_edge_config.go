package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
)

var save bool

var config gconfig.Config

func init() {
	flag.BoolVar(&save, "save", false, "persist config globally")
	config = gconfig.SourceAndParse()
}

func main() {

	jbyt, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(jbyt))

	if save {
		gconfig.Persist(config)
	}
}
