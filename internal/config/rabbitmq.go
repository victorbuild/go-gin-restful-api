package config

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp091.Connection
var RabbitChannel *amqp091.Channel

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
	RabbitConn, err = amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("❌ RabbitMQ 連線失敗:", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatal("❌ 開啟 RabbitMQ Channel 失敗:", err)
	}

	// 確保 `user_created` 佇列存在
	_, err = RabbitChannel.QueueDeclare(
		"user_created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ 無法建立 Queue:", err)
	}

	log.Println("✅ RabbitMQ 連線成功")
}

// publishMessageFunc 用於測試時替換的函數變數
var publishMessageFunc = func(queueName string, message string) {
	if RabbitChannel == nil {
		log.Println("⚠️ RabbitMQ Channel 未初始化，跳過訊息發送")
		return
	}
	err := RabbitChannel.PublishWithContext(
		nil,
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println("❌ 發送訊息失敗:", err)
		return
	}

	log.Println("✅ 訊息發送成功:", message)
}

// PublishMessage 發送訊息
func PublishMessage(queueName string, message string) {
	publishMessageFunc(queueName, message)
}

// SetPublishMessageFunc 設置 PublishMessage 的實現（測試用）
func SetPublishMessageFunc(fn func(string, string)) {
	publishMessageFunc = fn
}

// CloseRabbitMQ 關閉 RabbitMQ 連線
func CloseRabbitMQ() {
	if RabbitChannel != nil {
		RabbitChannel.Close()
	}
	if RabbitConn != nil {
		RabbitConn.Close()
	}
}
