package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig 設定 JWT Secret & Token 期限
var JWTConfig = struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}{
	AccessTokenSecret:  GetEnv("ACCESS_TOKEN_SECRET", "your-access-secret"),
	RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", "your-refresh-secret"),
	AccessTokenExpiry:  time.Second * time.Duration(GetEnvInt("ACCESS_TOKEN_EXPIRY", 3600)),    // 1 小時
	RefreshTokenExpiry: time.Second * time.Duration(GetEnvInt("REFRESH_TOKEN_EXPIRY", 604800)), // 7 天
}

// GenerateTokens 產生 Access Token & Refresh Token
func GenerateTokens(userID uint, email, role string) (string, string, error) {
	// 建立 Access Token
	accessClaims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(JWTConfig.AccessTokenExpiry).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(JWTConfig.AccessTokenSecret))
	if err != nil {
		return "", "", err
	}

	// 建立 Refresh Token
	refreshClaims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(JWTConfig.RefreshTokenExpiry).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(JWTConfig.RefreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateToken 驗證 Token 是否有效
func ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}
