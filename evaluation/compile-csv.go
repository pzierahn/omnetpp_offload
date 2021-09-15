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

	outdir := "data"
	_ = os.RemoveAll(outdir)
	_ = os.MkdirAll(outdir, 0755)

	csvs := []combine{
		{
			file:   "overhead.csv",
			prefix: "overhead",
			src:    "meta",
		},
		{
			file:   "setup.csv",
			prefix: "setup",
			src:    "meta",
		},
		{
			file:   "actions.csv",
			prefix: "actions",
			src:    "meta",
		},
		{
			file:   "transfers.csv",
			prefix: "transfers",
			src:    "meta",
		},
		{
			file:   "runs.csv",
			prefix: "runs",
			src:    "meta",
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
