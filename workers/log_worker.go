package workers

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"restfulapi/config"
)

func StartLogWorker() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetEnv("KAFKA_BROKER", "kafka:9092"),
		"group.id":          "log-consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal("❌ 無法建立 Kafka Consumer:", err)
	}
	defer c.Close()

	topic := config.GetEnv("KAFKA_TOPIC", "log-topic")

	// **檢查 Topic 是否存在**
	topics, err := c.GetMetadata(nil, true, 5000)
	if err != nil {
		log.Fatal("❌ 無法取得 Kafka Topic 列表:", err)
	}
	if _, exists := topics.Topics[topic]; !exists {
		log.Fatalf("❌ Kafka Topic `%s` 不存在，請先建立!", topic)
	}

	// **訂閱 Topic**
	c.SubscribeTopics([]string{topic}, nil)

	log.Println("🚀 `log_worker` 開始監聽 Kafka `log-topic` ...")
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("📩 Kafka Log 訊息: %s\n", string(msg.Value))
		} else {
			fmt.Printf("❌ 錯誤: %v (%v)\n", err, msg)
		}
	}
}
