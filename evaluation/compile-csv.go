package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type combine struct {
	file   string
	prefix string
	suffix string
}

func main() {

	// scp -r "pzierahn@85.214.35.83:/home/pzierahn/.cache/omnetpp-edge/evaluation/*.csv" meta

	outdir := "data"
	_ = os.RemoveAll(outdir)
	_ = os.MkdirAll(outdir, 0755)

	csvs := []combine{
		{
			file:   "durations.csv",
			prefix: "durations",
		},
		{
			file:   "setup.csv",
			suffix: "setup.csv",
		},
		{
			file:   "actions.csv",
			suffix: "actions.csv",
		},
		{
			file:   "transfers.csv",
			suffix: "transfers.csv",
		},
		{
			file:   "runs.csv",
			suffix: "runs.csv",
		},
	}

	for _, obj := range csvs {

		var entries []string
		var inx int

		err := filepath.WalkDir("meta", func(path string, file fs.DirEntry, _ error) (err error) {

			if !strings.HasPrefix(file.Name(), obj.prefix) {
				return
			}

			if !strings.HasSuffix(file.Name(), obj.suffix) {
				return
			}

			byt, err := os.ReadFile(path)
			if err != nil {
				return
			}

			if len(byt) == 0 {
				return
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

			return
		})
		if err != nil {
			panic(err)
		}

		data := strings.Join(entries, "\n") + "\n"
		err = ioutil.WriteFile(filepath.Join(outdir, obj.file), []byte(data), 0755)
		if err != nil {
			panic(err)
		}
	}
}
