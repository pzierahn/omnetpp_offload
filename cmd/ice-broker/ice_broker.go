package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var localHTTPPort = 5088

type ICEInfo struct {
	LocalUfrag    string
	LocalPwd      string
	Candidates    []string
	IsControlling bool
}

type Controlling struct {
	IsControlling bool
}

var controlling ICEInfo
var notControlling ICEInfo

func candidate(w http.ResponseWriter, req *http.Request) {

	var value ICEInfo
	err := json.NewDecoder(req.Body).Decode(&value)
	if err != nil {
		panic(err)
	}

	fmt.Println("candidate", "value", value)

	if value.IsControlling {
		controlling = value
	} else {
		notControlling = value
	}
}

func exchange(w http.ResponseWriter, req *http.Request) {

	var cont Controlling
	err := json.NewDecoder(req.Body).Decode(&cont)
	if err != nil {
		panic(err)
	}

	fmt.Println("exchange", "controlling", cont.IsControlling)

	var jbyt []byte

	if cont.IsControlling {
		jbyt, _ = json.Marshal(notControlling)
	} else {
		jbyt, _ = json.Marshal(controlling)
	}

	_, err = w.Write(jbyt)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Start broker")

	http.HandleFunc("/candidate", candidate)
	http.HandleFunc("/exchange", exchange)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", localHTTPPort), nil); err != nil {
		panic(err)
	}
}
