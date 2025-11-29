package config

// DatabaseConfig 設定資料庫變數
var DatabaseConfig = struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}{
	Host:     GetEnv("DB_HOST", "localhost"),
	Port:     GetEnv("DB_PORT", "5432"),
	User:     GetEnv("DB_USER", "postgres"),
	Password: GetEnv("DB_PASSWORD", "admin"),
	DBName:   GetEnv("DB_NAME", "postgres"),
}
