package kafka

import (
	"net"
	"os"
	"storage/config"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type CacheInterface interface {
	Add(string)
}

var kache CacheInterface
var kafkaURL string

// Создаем подключение к кафке, настраиваем ее, чекаем подключение
func Init(k CacheInterface) {

	kache = k

	kafkaURL = os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost"
	}

	conn, err := kafka.Dial("tcp", kafkaURL+":9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: config.GetString("topic-name"), NumPartitions: 1, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
