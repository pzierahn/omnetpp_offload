package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	start := ""
	out := "system-overhead.csv"
	dir := "system-overhead-scenarios"

	//dir := "system-overhead-data"

	//start := "opp-edge-eval-setup"
	//out := "opp-edge-eval-setup.csv"

	//start := "opp-edge-eval-transfers"
	//out := "opp-edge-eval-transfers.csv"

	//start := "opp-edge-eval-runs"
	//out := "opp-edge-eval-runs.csv"

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var entries []string
	var inx int

	for _, file := range files {

		if !strings.HasPrefix(file.Name(), start) {
			continue
		}

		byt, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			panic(err)
		}

		txt := string(byt)
		txt = strings.TrimSpace(txt)
		lines := strings.Split(txt, "\n")

		if inx == 0 {
			entries = append(entries, lines...)
		} else {
			entries = append(entries, lines[1:]...)
		}

		inx++
	}

	data := strings.Join(entries, "\n") + "\n"
	err = ioutil.WriteFile(out, []byte(data), 0755)
	if err != nil {
		panic(err)
	}
}
