package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims определяет структуру claims для JWT токенов
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey           string
	accessTokenDuration time.Duration
}

// NewJWTManager создает новый менеджер JWT токенов
func NewJWTManager(secretKey string, accessTokenHours int) *JWTManager {
	return &JWTManager{
		secretKey:           secretKey,
		accessTokenDuration: time.Duration(accessTokenHours) * time.Hour,
	}
}

// GenerateToken создает новый JWT токен для пользователя
func (manager *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(manager.accessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "ravr-backend",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// ValidateToken проверяет JWT токен и возвращает claims
func (manager *JWTManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
