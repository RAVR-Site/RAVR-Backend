package controller

import (
	"net/http"
	"strconv"

	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ResultController struct {
	svc    service.ResultService
	logger *zap.Logger
}

func NewResultController(svc service.ResultService, logger *zap.Logger) *ResultController {
	return &ResultController{
		svc:    svc,
		logger: logger,
	}
}

type saveResultRequest struct {
	TimeTaken uint   `json:"time_taken" example:"120"`
	LessonId  string `json:"lesson_id" example:"1"`
}

// @Description Ответ с данными одного урока
type SwaggerResultSaveResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    *service.Leaderboard `json:"data"`
}

// Save godoc
// @Summary      Save result
// @Description  Save result for a lesson
// @Tags         Results
// @Accept       json
// @Produce      json
// @Param        request body saveResultRequest true "Save result request"
// @Success      200 {object} SwaggerResultSaveResponse "Success"
// @Failure      400 {object} responses.ErrorResponse "Invalid request"
// @Failure      500 {object} responses.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/results/save [post]
// Save handles saving the result for a lesson.
func (s ResultController) Save(c echo.Context) error {
	var req saveResultRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error()))
	}

	username := c.Get("username").(string)

	// Проверяем, что lessonId может быть преобразован в uint (для GetLeaderboard)
	lessonIdUint, err := strconv.ParseUint(req.LessonId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", "LessonId must be an integer"))
	}

	// Вызываем Save с строковым lessonId
	err = s.svc.Save(username, req.LessonId, req.TimeTaken)
	if err != nil {
		s.logger.Error("Failed to save result", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, responses.Error("INTERNAL_ERROR", "Failed to save result"))
	}

	// Вызываем GetLeaderboard с числовым lessonId
	leaderboard, err := s.svc.GetLeaderboard(username, uint(lessonIdUint), 10)
	if err != nil {
		s.logger.Error("Failed to get leaderboard", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, responses.Error("INTERNAL_ERROR", "Failed to get leaderboard"))
	}

	return c.JSON(http.StatusOK, responses.Success(leaderboard))
}
