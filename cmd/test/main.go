package main

import (
	"github.com/patrickz98/project.go.omnetpp/simple"
)

const input = "/Users/patrick/Desktop/xxx/tictoc"
const mirror = "/Users/patrick/Desktop/xxx/tictoc-mirror"

func main() {

	//_ = os.RemoveAll(mirror)
	//_ = os.MkdirAll(mirror, 0755)

	_ = simple.SymbolicCopy(input, mirror+"-1", map[string]bool{
		"results/": true,
	})

	//file, err := os.Open("data/storage/tictoc-1f779d/source.tar.gz")
	//if err != nil {
	//	panic(err)
	//}
	//
	//_ = os.RemoveAll("data/xxx")
	//_ = os.MkdirAll("data/xxx", 0755)
	//
	//err = simple.UnTarGz("data/xxx", file)
	//if err != nil {
	//	panic(err)
	//}
}
