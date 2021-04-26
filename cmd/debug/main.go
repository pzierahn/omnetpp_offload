package main

import (
	"github.com/patrickz98/project.go.omnetpp/simple"
)

const input = "/Users/patrick/Desktop/xxx/tictoc"
const mirror = "/Users/patrick/Desktop/xxx/tictoc-mirror"

func main() {

	_, err := simple.TarGz("../TaskletSimulator", "test",
		"\\.git", ".DS_Store", ".idea", "out/clang-release/src/messages/DataRequestMessage_m.o")
	if err != nil {
		panic(err)
	}

	//path := "src/libTaskletSimulator.dylib"
	//
	//ext := filepath.Ext(path)
	//fmt.Println("ext", ext)
	//fmt.Println("ext", strings.TrimSuffix(path, ext))
	//fmt.Println("ext", strings.TrimSuffix("hallooo", "o"))
	//fmt.Println("Base", filepath.Base("/Users/patrick/github/project.go.omnetpp"))

	//fmt.Println("Base", filepath.Base("/Users/patrick/github/project.go.omnetpp"))

	//_ = os.RemoveAll(mirror)
	//_ = os.MkdirAll(mirror, 0755)

	//_ = simple.SymbolicCopy(input, mirror+"-1", map[string]bool{
	//	"results/": true,
	//})

	//file, err := os.Open("data/storage/tictoc-1f779d/source.tar.gz")
	//if err != nil {
	//	panic(err)
	//}

	//_ = os.RemoveAll("data/xxx")
	//_ = os.MkdirAll("data/xxx", 0755)

	//err = ioutil.WriteFile("data/xxx/xxx.tgz", buf.Bytes(), 0644)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = simple.UnTarGz("data/xxx", &buf)
	//if err != nil {
	//	panic(err)
	//}
}
