package config

// SMTP 設定
var (
	SMTPHost     = GetEnv("SMTP_HOST", "smtp.gmail.com")
	SMTPPort     = GetEnv("SMTP_PORT", "587")
	SMTPUser     = GetEnv("SMTP_USER", "your-email@gmail.com")
	SMTPPassword = GetEnv("SMTP_PASSWORD", "your-app-password")
)
