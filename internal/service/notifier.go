package service

import (
	"log"
	"net/smtp"
	"restfulapi/internal/config"
)

// Notifier 定義通知服務的介面，支援多種通知方式
type Notifier interface {
	Send(to string, subject string, body string) error
}

type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (n *EmailNotifier) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.SMTPUser, []string{to}, msg)
	if err != nil {
		log.Println("發送 Email 失敗:", err)
		return err
	}

	log.Println("Email 已發送到:", to)
	return nil
}

// SMSNotifier 簡訊通知實作（展示用）
// TODO: 未來實作時需要串接簡訊服務商 API
// 目前僅為展示介面設計模式，實際專案中若未使用可先移除
type SMSNotifier struct {
}

func NewSMSNotifier() *SMSNotifier {
	return &SMSNotifier{}
}

func (n *SMSNotifier) Send(to string, subject string, body string) error {
	log.Printf("發送簡訊到 %s: %s", to, body)
	return nil
}
