package main

import (
	"fmt"
	"net/http"
	"storage/db"
	"storage/kache"
	"storage/kafka"
	"time"
)

var kacheData *kache.StringSet

func hello(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprintf(w, "hello\n")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only GET equest"))
	}

}

func send(w http.ResponseWriter, req *http.Request) {
	kafka.SendMessage(time.Now().String())
	w.WriteHeader(http.StatusOK)

}

func allRecords(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	// Записываем каждую строку с новой строки
	for _, record := range kacheData.GetAll() {
		w.Write([]byte(record + "\n"))
	}

}

func main() {

	kacheData = kache.NewStringSet(1000)

	db.Init()
	kafka.Init()
	kafka.ReadKafka(kacheData)

	http.HandleFunc("/send", send)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/all", allRecords)

	http.ListenAndServe(":8080", nil)
}
