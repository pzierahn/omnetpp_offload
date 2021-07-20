package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/defines"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
)

var save bool
var showPaths bool

func init() {
	flag.BoolVar(&save, "save", false, "persist config globally")
	flag.BoolVar(&showPaths, "paths", false, "print paths")
}

func main() {

	gconfig.ParseFlags()

	if showPaths {
		fmt.Println("CacheDir:  ", defines.CacheDir())
		fmt.Println("ConfigDir: ", defines.ConfigDir())
		return
	}

	jbyt, _ := json.MarshalIndent(gconfig.Config, "", "  ")
	fmt.Println(string(jbyt))

	if save {
		gconfig.Persist()
	}
}
