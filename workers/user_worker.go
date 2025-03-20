package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"restfulapi/config"
	"restfulapi/models"
)

// 發送 Email（直接在 Worker 裡發送）
func sendEmail(to string, subject string, body string) {
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.SMTPUser, []string{to}, msg)
	if err != nil {
		log.Println("❌ 發送 Email 失敗:", err)
	} else {
		log.Println("✅ Email 已發送到:", to)
	}
}

// StartUserWorker 監聽 RabbitMQ
func StartUserWorker() {
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
		var user models.User
		err := json.Unmarshal(msg.Body, &user)
		if err != nil {
			log.Println("❌ JSON 解析失敗:", err)
			continue
		}

		// **發送 Email**
		subject := "歡迎加入我們！"
		body := fmt.Sprintf("嗨 %s，感謝你的註冊！", user.Name)
		sendEmail(user.Email, subject, body)

		fmt.Printf("📧 已發送 Email 給: %s (%s)\n", user.Name, user.Email)
	}
}
