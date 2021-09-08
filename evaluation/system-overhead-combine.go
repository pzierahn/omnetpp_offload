package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	//dir := "system-overhead-scenarios"
	//combine := map[string]string{
	//	"system-overhead.csv": "",
	//}

	dir := "system-overhead-data"
	combine := map[string]string{
		//"system-overhead.csv":         "",
		"opp-edge-eval-setup.csv":     "opp-edge-eval-setup",
		"opp-edge-eval-transfers.csv": "opp-edge-eval-transfers",
		"opp-edge-eval-runs.csv":      "opp-edge-eval-runs",
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for csvname, start := range combine {
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
		err = ioutil.WriteFile(csvname, []byte(data), 0755)
		if err != nil {
			panic(err)
		}
	}
}
