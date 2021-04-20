package main

import (
	"archive/tar"
	"bytes"
	"com.github.patrickz98.omnet/utils"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var logger *log.Logger

const (
	omnetBin = "/Users/patrick/Desktop/omnetpp-5.6.2/bin"
	tictoc   = "/Users/patrick/github/tictoc"
)

func init() {
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
}

func omnet() {

	simulationExe := "simulation"
	config := "TicToc18"

	//
	// Create Makefile
	//

	makemake := exec.Command(omnetBin+"/opp_makemake", "-f", "--deep", "-u", "Cmdenv", "-o", simulationExe)
	makemake.Dir = tictoc
	makemake.Stdout = os.Stdout
	makemake.Stderr = os.Stderr

	if err := makemake.Run(); err != nil {
		logger.Fatalln(err)
	}

	//
	// Compile simulation
	//

	makeCmd := exec.Command("sh", "-c", "make")
	makeCmd.Dir = tictoc
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	if err := makeCmd.Run(); err != nil {
		logger.Fatalln(err)
	}

	//
	// Get runnumbers
	//

	runnumbers := exec.Command("./"+simulationExe, "-c", config, "-s", "-q", "runnumbers")
	runnumbers.Dir = tictoc
	//runnumbers.Stdout = os.Stdout
	//runnumbers.Stderr = os.Stderr

	byt, err := runnumbers.CombinedOutput()
	if err != nil {
		logger.Fatalln(err)
	}

	output := string(byt)
	output = strings.TrimSpace(output)

	configs := strings.Split(output, " ")
	log.Println("configs", configs)

	//
	// Get configs
	//

	simConfigs := exec.Command("./"+simulationExe, "-c", config, "-s", "-a")
	simConfigs.Dir = tictoc
	//simConfigs.Stdout = os.Stdout
	//simConfigs.Stderr = os.Stderr

	byt, err = simConfigs.CombinedOutput()
	if err != nil {
		logger.Fatalln(err)
	}

	output = string(byt)
	output = strings.TrimSpace(output)

	reg := regexp.MustCompile(`Config (.+?):`)
	matches := reg.FindAllStringSubmatch(output, -1)

	configurations := make([]string, 0)

	for _, match := range matches {
		configurations = append(configurations, match[1])
	}

	logger.Println("matches", configurations)

	//
	// Run simulation
	//

	sim := exec.Command("./"+simulationExe, "-c", config, "-r", "1")
	sim.Dir = tictoc
	sim.Stdout = os.Stdout
	sim.Stderr = os.Stderr

	if err = sim.Run(); err != nil {
		logger.Fatalln(err)
	}
}

type Ping struct {
	DeviceName string
}

func ping(wri http.ResponseWriter, req *http.Request) {
	//_, err := fmt.Fprintf(wri, "hello\n")
	//if err != nil {
	//	logger.Println(err)
	//}

	logger.Println("req.URL.Path", req.URL.Path)
	logger.Println("req.Proto", req.Proto)
	logger.Println("req.RemoteAddr", req.RemoteAddr)

	http.ServeFile(wri, req, "go.mod")
}

func ZipWriter(id string) {
	var buffer bytes.Buffer

	zr := gzip.NewWriter(&buffer)
	tw := tar.NewWriter(zr)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			logger.Printf("Skip: %#v\n", path)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		logger.Printf("Crawling: %#v\n", path)

		input, err := os.Open(path)
		if err != nil {
			return err
		}

		// generate tar header
		header := &tar.Header{
			Name:    "tictoc-" + id + "/" + strings.TrimPrefix(path, tictoc),
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
		}

		// write header
		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		_, err = io.Copy(tw, input)
		if err != nil {
			return err
		}

		_ = input.Close()

		return nil
	}

	err := filepath.Walk(tictoc, walker)
	if err != nil {
		logger.Fatalln(err)
	}

	_ = tw.Close()
	_ = zr.Close()

	err = ioutil.WriteFile("tictoc-"+id+".tar.gz", buffer.Bytes(), 0766)
	if err != nil {
		logger.Fatalln(err)
	}
}

type Run struct {
	RunNumber string
	Status    string
}

type Storage struct {
	Id      string
	Created time.Time
	Source  string
	Configs map[string][]Run
}

func storage(wri http.ResponseWriter, req *http.Request) {

	id := utils.QueryString(req.URL.Query(), "id", "")

	if id == "" {
		wri.WriteHeader(400)
		return
	}

	logger.Println(req.Method, "id", id)

	if req.Method == "GET" {
		wri.Header().Add("Content-Type", "application/json")

		storage := Storage{
			Id:      "tictoc-e928",
			Created: time.Now(),
			Source:  "storage/tictoc-e928/source.tar.gz",
			Configs: map[string][]Run{
				"TicToc18": {
					{
						RunNumber: "1",
						Status:    "done",
					},
					{
						RunNumber: "2",
						Status:    "doing",
					},
					{
						RunNumber: "2",
						Status:    "todo",
					},
				},
			},
		}

		byt, err := json.MarshalIndent(storage, "", "  ")
		if err != nil {
			logger.Println(err)
			wri.WriteHeader(400)

			return
		}

		wri.WriteHeader(200)
		_, _ = wri.Write(byt)

		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	fileBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Println(err)
	}

	logger.Printf("fileBytes: %+v\n", len(fileBytes))
}

func uploadFile() {
	time.Sleep(time.Second * 2)

	filename := "/Users/patrick/Desktop/omnet/project.go.omnet/tictoc-e928.tar.gz"
	logger.Println("upload", filename)

	file, err := os.Open(filename)
	if err != nil {
		logger.Println(err)
	}

	_, err = http.Post("http://localhost:8090/storage?id=test-123", "application/x-gzip", file)
	if err != nil {
		logger.Println(err)
	}
}

func main() {

	//ZipWriter(simple.RandomId(4))

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/storage", storage)
	http.HandleFunc("/storage/*", ping)

	go uploadFile()

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		logger.Println(err)
	}
}
