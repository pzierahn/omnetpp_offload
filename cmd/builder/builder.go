package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	bpath = "build"
)

var build = map[string][]string{
	//"darwin": {
	//	"amd64",
	//	"arm64",
	//},
	"linux": {
		"amd64",
		"arm64",
	},
	//"windows": {
	//	"amd64",
	//	"386",
	//},
}

var cmds = map[string]string{
	"stargate":        "cmd/stargate/stargate.go",
	//"opp_edge_run":    "cmd/consumer/opp_edge_run.go",
	//"opp_edge_worker": "cmd/worker/opp_edge_worker.go",
	//"opp_edge_broker": "cmd/broker/opp_edge_broker.go",
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := os.RemoveAll(bpath); err != nil {
		log.Fatalln(err)
	}

	if err := os.MkdirAll(bpath, 0755); err != nil {
		log.Fatalln(err)
	}

	// go build -o build/caddy-server github.com/mholt/caddy/caddy
	// env GOOS=target-OS GOARCH=target-architecture go build package-import-path

	for cmd, path := range cmds {
		for tos, archs := range build {
			for _, arch := range archs {
				name := strings.Join([]string{cmd, tos, arch}, "_")
				obj := filepath.Join(bpath, name)

				log.Printf("building %v", obj)

				cmd := exec.Command("go", "build", "-o", obj, path)
				cmd.Env = os.Environ()
				cmd.Env = append(cmd.Env, "GOOS="+tos)
				cmd.Env = append(cmd.Env, "GOARCH="+arch)

				if byt, err := cmd.CombinedOutput(); err != nil {
					log.Println(string(byt))
					log.Fatalln(err)
				}
			}
		}
	}
}
