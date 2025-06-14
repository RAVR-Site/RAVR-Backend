package controller

import (
	"net/http"

	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CompleteRequest представляет запрос на завершение урока
type CompleteRequest struct {
	UserID           uint   `json:"userId"`
	LessonID         string `json:"lessonId"`
	CompletionTime   string `json:"completionTime"`
	EarnedExperience uint64 `json:"earnedExperience"`
}

// LessonCompletionController контроллер для работы с завершением уроков
type LessonCompletionController struct {
	svc    service.LessonCompletionService
	logger *zap.Logger
}

// NewLessonCompletionController создает новый контроллер для работы с завершением уроков
func NewLessonCompletionController(svc service.LessonCompletionService, logger *zap.Logger) *LessonCompletionController {
	return &LessonCompletionController{
		svc:    svc,
		logger: logger,
	}
}

// CompleteLessonResponse представляет ответ на запрос завершения урока
type CompleteLessonResponse struct {
	Success       bool   `json:"success"`
	Experience    uint64 `json:"experience"`    // Текущий опыт пользователя
	EarnedXP      uint64 `json:"earnedXP"`      // Заработанный опыт за текущий урок
	CompletedTime string `json:"completedTime"` // Время завершения урока
}

// Complete обрабатывает запрос на завершение урока
// @Summary Завершение урока
// @Description Обновляет статистику пользователя при завершении урока
// @Tags lessons
// @Accept json
// @Produce json
// @Param request body CompleteRequest true "Данные о прохождении урока"
// @Success 200 {object} CompleteLessonResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/lessons/complete [post]
func (c *LessonCompletionController) Complete(e echo.Context) error {
	var req CompleteRequest
	if err := e.Bind(&req); err != nil {
		c.logger.Error("Error binding request", zap.Error(err))
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error()))
	}

	result, err := c.svc.CompleteLessonAndUpdateExperience(req.UserID, req.LessonID, req.CompletionTime, req.EarnedExperience)
	if err != nil {
		c.logger.Error("Error completing lesson", zap.Error(err), zap.Uint("userId", req.UserID), zap.String("lessonId", req.LessonID))
		return e.JSON(http.StatusInternalServerError, responses.Error("COMPLETION_ERROR", err.Error()))
	}

	response := CompleteLessonResponse{
		Success:       true,
		Experience:    result.TotalExperience,
		EarnedXP:      result.EarnedExperience,
		CompletedTime: result.CompletionTime,
	}

	return e.JSON(http.StatusOK, responses.Success(response))
}
