package kafka

import (
	"context"
	"fmt"
	"storage/db"
	"time"

	"github.com/segmentio/kafka-go"
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
	r.SetOffsetAt(context.Background(), time.Now().Add(-24*time.Hour))
	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Error read message from kafka")
				break
			}
			db.WriteData(string(m.Value))
		}
	}()
}
