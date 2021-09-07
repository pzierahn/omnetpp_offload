package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	out := "system-overhead.csv"
	dir := "system-overhead-scenarios"

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var entries []string

	for inx, file := range files {
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
	}

	data := strings.Join(entries, "\n") + "\n"
	err = ioutil.WriteFile(out, []byte(data), 0755)
	if err != nil {
		panic(err)
	}
}
