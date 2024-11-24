package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Глобальная переменная для хранения продюсера Kafka
var (
	writer *kafka.Writer
	once   sync.Once
)

// GetKafkaWriter возвращает продюсера Kafka, создавая его один раз
func GetKafkaWriter(brokerURL, topic string) *kafka.Writer {
	once.Do(func() {
		writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{brokerURL},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
	})

	return writer
}

// SendMessage отправляет сообщение "привет" в указанный топик Kafka.
func sendMessage(brokerURL, topic string, message string) error {
	// Одно подключение на всё время жизни приложения

	writer = GetKafkaWriter(brokerURL, topic)

	// defer writer.Close()

	// Отправка сообщения
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

func SendMessage(message string) {
	// Получение URL брокера Kafka из переменной окружения или использование значения по умолчанию
	brokerURL := os.Getenv("KAFKA_URL")
	if brokerURL == "" {
		brokerURL = "localhost:9092"
	}

	// Указание топика
	topic := "storage"

	// Отправка сообщения
	err := sendMessage(brokerURL, topic, message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
