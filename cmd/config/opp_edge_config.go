package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
)

var save bool

func init() {
	flag.BoolVar(&save, "save", false, "persist config globally")
	gconfig.ParseFlags()
}

func main() {

	jbyt, _ := json.MarshalIndent(gconfig.Config, "", "  ")
	fmt.Println(string(jbyt))

	if save {
		gconfig.Persist()
	}
}
