package simple

import (
	"log"
	"net/http"
)

func Watch(path string, onRequest func() interface{}) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		byt := PrettyBytes(onRequest())
		_, err := w.Write(byt)
		if err != nil {
			log.Fatalln(err)
		}
	})
}

func StartWatchServer(addr string) {

	if addr == "" {
		addr = ":8077"
	}

	log.Printf("start watch server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
