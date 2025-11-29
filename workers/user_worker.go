package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"restfulapi/internal/config"
	"restfulapi/internal/model"
	"restfulapi/internal/service"
)

// StartUserWorker 監聽 RabbitMQ
func StartUserWorker() {
	emailNotifier := service.NewEmailNotifier()
	config.InitRabbitMQ()
	defer config.CloseRabbitMQ()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgs, err := config.RabbitChannel.ConsumeWithContext(
		ctx,
		"user_created",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ 訂閱 Queue 失敗:", err)
	}

	log.Println("🚀 `user_worker` 開始監聽 user_created 事件...")

	for msg := range msgs {
		var user model.User
		err := json.Unmarshal(msg.Body, &user)
		if err != nil {
			log.Println("❌ JSON 解析失敗:", err)
			continue
		}

		// **發送 Email**
		subject := "歡迎加入我們！"
		body := fmt.Sprintf("嗨 %s，感謝你的註冊！", user.Name)
		emailNotifier.Send(user.Email, subject, body)

		fmt.Printf("📧 已發送 Email 給: %s (%s)\n", user.Name, user.Email)
	}
}
