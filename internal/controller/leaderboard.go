package controller

import (
	"net/http"
	"strconv"

	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Типы для Swagger документации
// @Description Ответ с таблицей лидеров
type SwaggerLeaderboardResponse struct {
	Success bool `json:"success" example:"true"`
	Data    []struct {
		UserID     uint   `json:"user_id" example:"1"`
		Username   string `json:"username" example:"johndoe"`
		FirstName  string `json:"first_name,omitempty" example:"John"`
		LastName   string `json:"last_name,omitempty" example:"Doe"`
		Position   int    `json:"position" example:"1"`
		Experience uint64 `json:"experience" example:"1500"`
		Trend      string `json:"trend" example:"up"`
	} `json:"data"`
}

// @Description Успешный ответ с сообщением для Swagger
type SwaggerSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    map[string]string `json:"data"`
}

// LeaderboardController контроллер для работы с таблицей лидеров
type LeaderboardController struct {
	svc    service.LeaderboardService
	logger *zap.Logger
}

// NewLeaderboardController создает новый контроллер для работы с таблицей лидеров
func NewLeaderboardController(svc service.LeaderboardService, logger *zap.Logger) *LeaderboardController {
	return &LeaderboardController{
		svc:    svc,
		logger: logger,
	}
}

// GetLeaderboard возвращает список топ-N пользователей по опыту
// @Summary Таблица лидеров
// @Description Возвращает список пользователей, отсортированных по опыту, с изменением позиции
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param limit query int false "Максимальное количество записей (по умолчанию 10)" default(10)
// @Success 200 {object} SwaggerLeaderboardResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/leaderboard [get]
func (c *LeaderboardController) GetLeaderboard(e echo.Context) error {
	// Получаем параметр limit из запроса, по умолчанию 10
	limitStr := e.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Получаем данные таблицы лидеров
	entries, err := c.svc.GetLeaderboard(limit)
	if err != nil {
		c.logger.Error("Failed to get leaderboard", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.Error("SERVER_ERROR", err.Error()))
	}

	return e.JSON(http.StatusOK, responses.Success(entries))
}

// UpdateRankings обновляет рейтинги пользователей за указанный период
// @Summary Обновление рейтингов
// @Description Обновляет рейтинги пользователей за указанный период (weekly, monthly)
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param period query string false "Период (daily, weekly, monthly)" default(weekly)
// @Success 200 {object} SwaggerSuccessResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/admin/leaderboard/update [post]
func (c *LeaderboardController) UpdateRankings(e echo.Context) error {
	// Получаем параметр period из запроса, по умолчанию weekly
	period := e.QueryParam("period")
	if period == "" {
		period = "weekly"
	}

	// Проверяем, что period имеет допустимое значение
	if period != "daily" && period != "weekly" && period != "monthly" {
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_PARAMETER", "Invalid period value"))
	}

	// Обновляем рейтинги
	err := c.svc.UpdateUserRankings(period)
	if err != nil {
		c.logger.Error("Failed to update rankings", zap.Error(err), zap.String("period", period))
		return e.JSON(http.StatusInternalServerError, responses.Error("SERVER_ERROR", err.Error()))
	}

	return e.JSON(http.StatusOK, responses.MessageResponse("Rankings updated successfully"))
}
