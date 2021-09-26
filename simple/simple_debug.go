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

func StartWatchServer() {
	log.Fatal(http.ListenAndServe(":8077", nil))
}
