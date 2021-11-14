package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/gconfig"
)

var save bool
var showPaths bool

func init() {
	flag.BoolVar(&save, "save", false, "persist config globally")
	flag.BoolVar(&showPaths, "paths", false, "print paths")
}

func main() {

	config := gconfig.ParseFlags()

	if showPaths {
		fmt.Println("CacheDir:  ", gconfig.CacheDir())
		fmt.Println("ConfigDir: ", gconfig.ConfigDir())
		return
	}

	jbyt, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(jbyt))

	if save {
		gconfig.Write()
	}
}
