package middleware

import (
	"net/http"
	"strings"

	"github.com/Ravr-Site/Ravr-Backend/internal/auth"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// JWTMiddleware создает middleware для проверки JWT токенов
func JWTMiddleware(jwtSecret string, jwtAccessExpiration int, logger *zap.Logger) echo.MiddlewareFunc {
	jwtManager := auth.NewJWTManager(jwtSecret, jwtAccessExpiration)
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			
			if authHeader == "" {
				logger.Debug("missing authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Проверяем формат "Bearer <token>"
			if !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Debug("invalid authorization header format")
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				logger.Debug("empty token")
				return echo.NewHTTPError(http.StatusUnauthorized, "empty token")
			}

			// Валидируем токен с помощью JWT менеджера
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				logger.Debug("invalid token", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Добавляем данные пользователя в контекст
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			logger.Debug("token validated successfully", 
				zap.String("username", claims.Username),
				zap.Uint("user_id", claims.UserID))

			return next(c)
		}
	}
}
