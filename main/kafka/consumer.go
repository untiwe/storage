package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"storage/cache"
	"storage/config"
	"storage/conventions"
	"storage/db"
)

// Запускаем ридер для кафки
func ReadKafka() {
	// Создание ридера Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL + ":9092"},
		Topic:     config.GetString("topic-name"),
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	//Не читаем предыдыущие сообщения. Они должны быть в базе.
	r.SetOffsetAt(context.Background(), time.Now())

	//Читаем бесконечно в фоне
	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("Error read message from kafka", err)
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
			cache.Add(string(m.Value))
		}
	}()
}
