package main

import (
	"fmt"
	"net/http"
	"storage/db"
	"storage/kafka"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprintf(w, "hello\n")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only GET equest"))
	}

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func send(w http.ResponseWriter, req *http.Request) {
	kafka.SendMessage(time.Now().String())
	w.WriteHeader(http.StatusOK)

}

func main() {

	db.Init()
	kafka.Init()
	kafka.ReadKafka()

	http.HandleFunc("/send", send)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}
