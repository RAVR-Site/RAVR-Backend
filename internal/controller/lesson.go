package controller

import (
	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// LessonController контроллер для работы с уроками
type LessonController struct {
	svc    service.LessonService
	logger *zap.Logger
}

// NewLessonController создает новый контроллер для уроков
func NewLessonController(svc service.LessonService, logger *zap.Logger) *LessonController {
	return &LessonController{
		svc:    svc,
		logger: logger,
	}
}

// @Description Ответ с данными одного урока
type SwaggerLessonResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    *service.Lesson `json:"data"`
}

// GetLesson возвращает детали конкретного урока
// @Summary Получение детальной информации урока
// @Description Возвращает детальную информацию об уроке по его ID, включая полные данные урока
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "ID урока"
// @Success 200 {object} SwaggerLessonResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /lessons/{id} [get]
func (c *LessonController) GetLesson(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "Некорректный ID урока"))
	}

	lesson, err := c.svc.GetLesson(uint(id))
	if err != nil {
		c.logger.Error("Ошибка получения урока", zap.String("id", idStr), zap.Error(err))
		return e.JSON(http.StatusNotFound, responses.Error("LESSON_NOT_FOUND", "Урок не найден"))
	}

	return e.JSON(http.StatusOK, responses.Success(lesson))
}

// @Description Ответ со списком уроков
type SwaggerLessonsResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    []*service.Lesson `json:"data"`
}

// GetLessonsByType возвращает уроки по типу
// @Summary Получение уроков по типу
// @Description Возвращает список уроков определенного типа
// @Tags lessons
// @Accept json
// @Produce json
// @Param type query string true "Тип урока"
// @Success 200 {object} SwaggerLessonsResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /lessons [get]
func (c *LessonController) GetLessonsByType(e echo.Context) error {
	lessonType := e.QueryParam("type")

	if lessonType == "" {
		return e.JSON(http.StatusBadRequest, responses.Error("INVALID_QUERY", "Параметры type и mode обязательны"))
	}

	lessons, err := c.svc.GetLessonByType(lessonType)
	if err != nil {
		c.logger.Error("Ошибка получения уроков", zap.String("type", lessonType), zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.Error("INTERNAL_ERROR", "Ошибка получения уроков"))
	}

	return e.JSON(http.StatusOK, responses.Success(lessons))
}
