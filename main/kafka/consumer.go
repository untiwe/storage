package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	// "log"

	"storage/conventions"
	"storage/db"
	"time"

	"github.com/segmentio/kafka-go"
	// "golang.org/x/text/message"
)

// https://github.com/segmentio/kafka-go
func ReadKafka() {
	// Создание ридера Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "storage",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	//Не читаем предыдыущие сообщения. Они должны быть в базе.
	r.SetOffsetAt(context.Background(), time.Now())
	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Error read message from kafka")
				break
			}

			var order conventions.Order
			err = json.Unmarshal([]byte(m.Value), &order)
			if err != nil {
				fmt.Println("Error unmarshal message:", err)
			}

			//пробуем записать в базу
			err = db.InsertOrder(order)
			if err != nil {
				fmt.Println("Order is not Insert:", err)
				return
			}

			//если удачно, дополняем кеш
			kache.Add(string(m.Value))
		}
	}()
}
