package controller

import (
	"net/http"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CompleteRequest представляет запрос на завершение урока
type CompleteRequest struct {
	LessonID         string `json:"lessonId"`
	CompletionTime   uint64 `json:"completionTime"` // Время в секундах
	EarnedExperience uint64 `json:"earnedExperience"`
}

// LessonCompletionController контроллер для работы с завершением уроков
type LessonCompletionController struct {
	svc            service.LessonCompletionService
	leaderboardSvc service.LeaderboardService
	logger         *zap.Logger
}

// NewLessonCompletionController создает новый контроллер для работы с завершением уроков
func NewLessonCompletionController(
	svc service.LessonCompletionService,
	leaderboardSvc service.LeaderboardService,
	logger *zap.Logger,
) *LessonCompletionController {
	return &LessonCompletionController{
		svc:            svc,
		leaderboardSvc: leaderboardSvc,
		logger:         logger,
	}
}

// CompleteLessonResponse представляет ответ на запрос завершения урока
type CompleteLessonResponse struct {
	Experience    uint64 `json:"experience"`    // Текущий опыт пользователя
	EarnedXP      uint64 `json:"earnedXP"`      // Заработанный опыт за текущий урок
	CompletedTime uint64 `json:"completedTime"` // Время завершения урока в секундах
	Leaderboard   struct {
		Entries      []*repository.LessonLeaderboardEntry `json:"entries"`      // Записи лидерборда
		UserPosition int                                  `json:"userPosition"` // Позиция пользователя
		TotalUsers   int                                  `json:"totalUsers"`   // Общее количество пользователей
	} `json:"leaderboard"` // Лидерборд по уроку
}

// Complete обрабатывает запрос на завершение урока
// @Summary Завершение урока
// @Description Обновляет статистику пользователя при завершении урока и возвращает лидерборд по уроку
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CompleteRequest true "Данные о прохождении урока"
// @Success 200 {object} CompleteLessonResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /lessons/complete [post]
func (c *LessonCompletionController) Complete(e echo.Context) error {
	var req CompleteRequest
	if err := e.Bind(&req); err != nil {
		c.logger.Error("Error binding request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error()))
	}

	// Получаем ID пользователя из контекста (установлено JWT middleware)
	userID := e.Get("user_id").(uint)

	result, err := c.svc.CompleteLessonAndUpdateExperience(userID, req.LessonID, req.CompletionTime, req.EarnedExperience)
	if err != nil {
		c.logger.Error("Error completing lesson", zap.Error(err), zap.Uint("userId", userID), zap.String("lessonId", req.LessonID))
		return e.JSON(http.StatusInternalServerError, responses.Error("COMPLETION_ERROR", err.Error()))
	}

	// Получаем данные для лидерборда
	limit := 10 // Показываем топ-10 в лидерборде
	entries, userPosition, err := c.leaderboardSvc.GetLessonLeaderboard(req.LessonID, userID, limit)
	if err != nil {
		c.logger.Error("Error getting leaderboard", zap.Error(err), zap.String("lessonId", req.LessonID))
		return e.JSON(http.StatusInternalServerError, responses.Error("LEADERBOARD_ERROR", err.Error()))
	}

	response := CompleteLessonResponse{
		Experience:    result.TotalExperience,
		EarnedXP:      result.EarnedExperience,
		CompletedTime: result.CompletionTime,
		Leaderboard: struct {
			Entries      []*repository.LessonLeaderboardEntry `json:"entries"`
			UserPosition int                                  `json:"userPosition"`
			TotalUsers   int                                  `json:"totalUsers"`
		}{
			Entries:      entries,
			UserPosition: userPosition,
			TotalUsers:   len(entries),
		},
	}

	return e.JSON(http.StatusOK, responses.Success(response))
}
