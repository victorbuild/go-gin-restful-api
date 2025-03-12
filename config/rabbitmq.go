package config

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ() {
	// **從 `.env` 讀取 RabbitMQ 設定**
	rabbitMQURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		GetEnv("RABBITMQ_USER", "admin"),
		GetEnv("RABBITMQ_PASSWORD", "admin"),
		GetEnv("RABBITMQ_HOST", "localhost"),
		GetEnv("RABBITMQ_PORT", "5672"),
		GetEnv("RABBITMQ_VHOST", "/"),
	)

	var err error
	RabbitConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("❌ RabbitMQ 連線失敗:", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatal("❌ 開啟 RabbitMQ Channel 失敗:", err)
	}

	// 確保 `user_created` 佇列存在
	_, err = RabbitChannel.QueueDeclare("user_created", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ 無法建立 Queue:", err)
	}

	log.Println("✅ RabbitMQ 連線成功")
}

// PublishMessage 發送訊息
func PublishMessage(queueName string, message string) {
	err := RabbitChannel.Publish(
		"", queueName, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(message)},
	)
	if err != nil {
		log.Println("❌ 發送訊息失敗:", err)
	} else {
		log.Println("✅ 訊息發送成功:", message)
	}
}
