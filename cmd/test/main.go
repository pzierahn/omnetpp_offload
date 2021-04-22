package main

import (
	"com.github.patrickz98.omnet/simple"
	"os"
)

func main() {
	file, err := os.Open("data/storage/tictoc-1f779d/source.tar.gz")
	if err != nil {
		panic(err)
	}

	_ = os.RemoveAll("data/xxx")
	_ = os.MkdirAll("data/xxx", 0755)

	err = simple.UnTarGz("data/xxx", file)
	if err != nil {
		panic(err)
	}
}
