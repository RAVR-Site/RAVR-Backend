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

// @Description Ответ с расширенной таблицей лидеров
type SwaggerExtendedLeaderboardResponse struct {
	Success bool `json:"success" example:"true"`
	Data    []struct {
		UserID           uint   `json:"user_id" example:"1"`
		Username         string `json:"username" example:"johndoe"`
		FirstName        string `json:"first_name,omitempty" example:"John"`
		LastName         string `json:"last_name,omitempty" example:"Doe"`
		Position         int    `json:"position" example:"1"`
		Experience       uint64 `json:"experience" example:"1500"`
		TotalLessons     int64  `json:"total_lessons" example:"15"`
		TotalTimeSpent   uint64 `json:"total_time_spent" example:"3540"`
		AverageTimeSpent string `json:"average_time_spent" example:"03:56"`
		Trend            string `json:"trend" example:"up"`
	} `json:"data"`
}

// @Description Ответ с таблицей лидеров конкретного урока
type SwaggerLessonLeaderboardResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		Entries []struct {
			UserID         uint   `json:"user_id" example:"1"`
			Username       string `json:"username" example:"johndoe"`
			FirstName      string `json:"first_name,omitempty" example:"John"`
			LastName       string `json:"last_name,omitempty" example:"Doe"`
			Position       int    `json:"position" example:"1"`
			CompletionTime uint64 `json:"completion_time" example:"85"` // Время прохождения урока в секундах
			Score          uint64 `json:"score" example:"85"`
			Experience     uint64 `json:"experience" example:"120"`
			Trend          string `json:"trend" example:"up"`
		} `json:"entries"`
		UserPosition int `json:"user_position" example:"3"`
		TotalUsers   int `json:"total_users" example:"50"`
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
// @Router /leaderboard [get]
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

// GetExtendedLeaderboard возвращает расширенную таблицу лидеров с информацией о времени
// @Summary Расширенная таблица лидеров
// @Description Возвращает расширенную таблицу лидеров с информацией о времени, затраченном на уроки
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param limit query int false "Максимальное количество записей (по умолчанию 10)" default(10)
// @Success 200 {object} SwaggerExtendedLeaderboardResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /leaderboard/extended [get]
func (c *LeaderboardController) GetExtendedLeaderboard(e echo.Context) error {
	// Получаем параметр limit из запроса, по умолчанию 10
	limitStr := e.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Получаем данные расширенной таблицы лидеров
	entries, err := c.svc.GetExtendedLeaderboard(limit)
	if err != nil {
		c.logger.Error("Failed to get extended leaderboard", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.Error("SERVER_ERROR", err.Error()))
	}

	return e.JSON(http.StatusOK, responses.Success(entries))
}

// GetLessonLeaderboard возвращает таблицу лидеров для конкретного урока
// @Summary Таблица лидеров урока
// @Description Возвращает таблицу лидеров для конкретного урока и позицию пользователя
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param lesson_id path string true "ID урока"
// @Param limit query int false "Максимальное количество записей (по умолчанию 10)" default(10)
// @Security BearerAuth
// @Success 200 {object} SwaggerLessonLeaderboardResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /leaderboard/lesson/{lesson_id} [get]
func (c *LeaderboardController) GetLessonLeaderboard(e echo.Context) error {
	// Получаем ID урока из пути
	lessonID := e.Param("lesson_id")
	if lessonID == "" {
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_PARAMETER", "Lesson ID is required"))
	}

	// Получаем ID пользователя из JWT-токена
	userID := e.Get("user_id").(uint)

	// Получаем параметр limit из запроса, по умолчанию 10
	limitStr := e.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Получаем таблицу лидеров урока
	entries, userPosition, err := c.svc.GetLessonLeaderboard(lessonID, userID, limit)
	if err != nil {
		c.logger.Error("Failed to get lesson leaderboard",
			zap.Error(err),
			zap.String("lessonID", lessonID),
			zap.Uint("userID", userID))
		return e.JSON(http.StatusInternalServerError, responses.Error("SERVER_ERROR", err.Error()))
	}

	// Формируем ответ
	response := struct {
		Entries      interface{} `json:"entries"`
		UserPosition int         `json:"user_position"`
		TotalUsers   int         `json:"total_users"` // Для простоты используем длину лидерборда, в реальности нужно считать
	}{
		Entries:      entries,
		UserPosition: userPosition,
		TotalUsers:   len(entries),
	}

	return e.JSON(http.StatusOK, responses.Success(response))
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
