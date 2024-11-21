package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// https://github.com/segmentio/kafka-go
func ReadKafka() {

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "storage", 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	b := make([]byte, 1e3) // 10KB max per message
	go func() {
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			batch := conn.ReadBatch(1e3, 1e6) // fetch 10KB min, 1MB max
			n, err := batch.Read(b)
			if err != nil {
				continue
			}
			fmt.Println(string(b[:n]))
		}
	}()

	// Создание потребителя
	// reader := kafka.NewReader(readerConfig)
	// defer reader.Close()

	// go func() {
	// 	for {
	// 		msg, err := reader.ReadMessage(context.Background())
	// 		if err != nil {
	// 			log.Printf("Could not read message: %v", err)
	// 			continue
	// 		}
	// 		fmt.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
	// 	}
	// }()
}
