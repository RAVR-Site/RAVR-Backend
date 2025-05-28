package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 1) // 1 час

	token, err := manager.GenerateToken(123, "testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 1) // 1 час

	// Генерируем токен
	token, err := manager.GenerateToken(123, "testuser")
	assert.NoError(t, err)

	// Валидируем токен
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "ravr-backend", claims.Issuer)
	assert.Equal(t, "testuser", claims.Subject)
}

func TestJWTManager_ValidateExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", -1) // -1 час (уже истек)

	// Генерируем токен, который уже истек
	token, err := manager.GenerateToken(123, "testuser")
	assert.NoError(t, err)

	// Ждем немного для гарантии истечения
	time.Sleep(time.Millisecond * 10)

	// Валидируем токен (должен быть недействительным)
	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTManager_ValidateTokenWithWrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret1", 1)
	manager2 := NewJWTManager("secret2", 1)

	// Генерируем токен с одним секретом
	token, err := manager1.GenerateToken(123, "testuser")
	assert.NoError(t, err)

	// Пытаемся валидировать с другим секретом
	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTManager_ValidateInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 1)

	// Пытаемся валидировать некорректный токен
	_, err := manager.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}
