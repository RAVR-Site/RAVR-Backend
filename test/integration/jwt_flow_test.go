package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/Ravr-Site/Ravr-Backend/internal/controller"
	"github.com/Ravr-Site/Ravr-Backend/internal/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() (*echo.Echo, error) {
	// Создаем тестовую конфигурацию
	cfg := &config.Config{
		JWTSecret:           "test-secret-for-jwt-testing",
		JWTAccessExpiration: 3600, // 1 час в секундах
	}

	// Создаем in-memory SQLite базу
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Мигрируем схему
	err = db.AutoMigrate(&repository.User{}, &repository.Lesson{})
	if err != nil {
		return nil, err
	}

	// Создаем logger
	logger := zap.NewNop()

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	
	// Создаем сервисы
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTAccessExpiration, logger)

	// Создаем контроллеры
	userController := controller.NewUserController(userService, logger)

	// Настраиваем Echo
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// Настраиваем роуты
	api := e.Group("/api/v1")
	
	// Публичные роуты
	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)
	
	// Защищенные роуты
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret, cfg.JWTAccessExpiration, logger))
	protected.GET("/user", userController.Profile)

	return e, nil
}

func TestJWTAuthFlow(t *testing.T) {
	e, err := setupTestApp()
	require.NoError(t, err)

	t.Run("Complete JWT Authentication Flow", func(t *testing.T) {
		// Шаг 1: Регистрируем нового пользователя
		registerPayload := map[string]interface{}{
			"username": "testuser",
			"password": "testpassword123",
		}
		registerBody, _ := json.Marshal(registerPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(registerBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		
		var registerResponse map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &registerResponse)
		require.NoError(t, err)
		
		// Проверяем формат ответа (может быть success: true вместо status: "success")
		if success, ok := registerResponse["success"].(bool); ok {
			assert.True(t, success)
		} else {
			assert.Equal(t, "success", registerResponse["status"])
		}

		// Шаг 2: Входим в систему и получаем JWT токен
		loginPayload := map[string]interface{}{
			"username": "testuser",
			"password": "testpassword123",
		}
		loginBody, _ := json.Marshal(loginPayload)

		req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var loginResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &loginResponse)
		require.NoError(t, err)
		
		// Проверяем формат ответа
		if success, ok := loginResponse["success"].(bool); ok {
			assert.True(t, success)
		} else {
			assert.Equal(t, "success", loginResponse["status"])
		}

		// Извлекаем JWT токен
		data, ok := loginResponse["data"].(map[string]interface{})
		require.True(t, ok)
		token, ok := data["token"].(string)
		require.True(t, ok)
		require.NotEmpty(t, token)

		fmt.Printf("Generated JWT Token: %s\n", token)

		// Шаг 3: Используем JWT токен для доступа к защищенному эндпоинту
		req = httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &profileResponse)
		require.NoError(t, err)
		
		// Проверяем формат ответа
		if success, ok := profileResponse["success"].(bool); ok {
			assert.True(t, success)
		} else {
			assert.Equal(t, "success", profileResponse["status"])
		}

		// Проверяем, что в ответе есть информация о пользователе
		profileData, ok := profileResponse["data"].(map[string]interface{})
		require.True(t, ok)
		username, ok := profileData["username"].(string)
		require.True(t, ok)
		assert.Equal(t, "testuser", username)

		fmt.Printf("Profile Response: %+v\n", profileResponse)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		// Тестируем с недействительным токеном
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Missing JWT Token", func(t *testing.T) {
		// Тестируем без токена
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// Тест истекшего токена требует сложной настройки времени
	// Основная функциональность JWT проверена в модульных тестах auth/jwt_test.go
}
