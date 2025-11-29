package config

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer Kafka Producer
var KafkaProducer *kafka.Producer

func InitKafka() {
	var err error
	KafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": GetEnv("KAFKA_BROKER", "localhost:9092"),
	})
	if err != nil {
		log.Fatal("❌ Kafka 連線失敗:", err)
	}

	log.Println("✅ Kafka 連線成功！")
}

// PublishKafkaMessage 發送 Kafka 訊息
func PublishKafkaMessage(topic string, message string) {
	err := KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Println("❌ 發送 Kafka 訊息失敗:", err)
	} else {
		log.Println("✅ Kafka 訊息已發送:", message)
	}
}
