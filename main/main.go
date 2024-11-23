package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storage/conventions"
	"storage/db"
	"storage/kache"
	"storage/kafka"
	// "time"
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

// func send(w http.ResponseWriter, req *http.Request) {
// 	kafka.SendMessage(time.Now().String())
// 	w.WriteHeader(http.StatusOK)

// }

func allRecords(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "aplication/json")

	recordsJSON, err := json.Marshal(kacheData.GetAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Записываем каждую строку с новой строки
	w.Write(recordsJSON)

}

func addOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}

	//базовая валидация
	var order conventions.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "JSON is not valid", http.StatusBadRequest)
		return
	}

	//если успешно, отправляем в брокер
	orderJSON, _ := json.Marshal(order)
	kafka.SendMessage(string(orderJSON))

	w.WriteHeader(http.StatusOK)

}

func main() {

	kacheData = kache.NewStringSet(1000)

	db.Init(kacheData)
	db.FillСache()
	kafka.Init(kacheData)
	kafka.ReadKafka()

	http.HandleFunc("/addorder", addOrder)
	// http.HandleFunc("/send", send)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/all", allRecords)

	http.ListenAndServe(":8080", nil)
}
