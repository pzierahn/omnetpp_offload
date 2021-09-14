package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type combine struct {
	file   string
	prefix string
	src    string
}

func main() {

	outdir := "eval-data"
	_ = os.RemoveAll(outdir)
	_ = os.MkdirAll(outdir, 0755)

	csvs := []combine{
		{
			file: "system-overhead.csv",
			src:  "system-overhead-scenarios",
		},
		{
			file:   "opp-edge-eval-setup.csv",
			prefix: "opp-edge-eval-setup",
			src:    "system-overhead-data",
		},
		{
			file:   "opp-edge-eval-actions.csv",
			prefix: "opp-edge-eval-actions",
			src:    "system-overhead-data",
		},
		{
			file:   "opp-edge-eval-transfers.csv",
			prefix: "opp-edge-eval-transfers",
			src:    "system-overhead-data",
		},
		{
			file:   "opp-edge-eval-runs.csv",
			prefix: "opp-edge-eval-runs",
			src:    "system-overhead-data",
		},
	}

	for _, obj := range csvs {

		files, err := os.ReadDir(obj.src)
		if err != nil {
			panic(err)
		}

		var entries []string
		var inx int

		for _, file := range files {

			if !strings.HasPrefix(file.Name(), obj.prefix) {
				continue
			}

			byt, err := os.ReadFile(filepath.Join(obj.src, file.Name()))
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
		err = ioutil.WriteFile(filepath.Join(outdir, obj.file), []byte(data), 0755)
		if err != nil {
			panic(err)
		}
	}
}
