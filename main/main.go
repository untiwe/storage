package main

import (
	"encoding/json"
	"net/http"

	"storage/cache"
	"storage/conventions"
	"storage/db"
	"storage/kafka"

)

// получить JSON со всеми записями
func allRecords(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "aplication/json")
	var records []conventions.Order
	for _, record := range cache.GetAll() {
		records = append(records, record)
	}

	recordsJSON, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Записываем каждую строку с новой строки
	w.Write(recordsJSON)

}

// Сгенерировать тестовый order и отправить в брокер
func generateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get", http.StatusMethodNotAllowed)
		return
	}
	order := conventions.GenerateRandomOrder()

	orderJSON, _ := json.Marshal(order)
	go func() {
		kafka.SendMessage(string(orderJSON))
	}()

	w.Write(orderJSON)

}

// Добавить order по полученому JSON
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

	//проверяем налчиче в кеше
	if cache.GetByID(order.OrderUID).OrderUID != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(order.OrderUID + " already used"))
		return
	}

	//если успешно, отправляем в брокер
	orderJSON, _ := json.Marshal(order)
	kafka.SendMessage(string(orderJSON))

	w.WriteHeader(http.StatusOK)

}

// вернуть странцу для отображения данных
func index(w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query()
	print(d)
	id := r.URL.Query().Get("id")
	if id != "" {
		order := cache.GetByID(id)
		if order.OrderUID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		orders := [1]conventions.Order{order}
		record, err := json.Marshal(orders)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.Write(record)
		return

	}

	http.ServeFile(w, r, "index.html")
}

func main() {

	//Явно вызываем инициализацию
	cache.Init()
	db.Init()
	kafka.Init()

	db.FillСache()
	kafka.ReadKafka()

	http.HandleFunc("/addorder", addOrder)
	http.HandleFunc("/generateorder", generateOrder)
	http.HandleFunc("/all", allRecords)
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}
