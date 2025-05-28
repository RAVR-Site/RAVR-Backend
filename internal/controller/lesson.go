package controller

import (
	"encoding/json"
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

// GetLessonTypes возвращает список уникальных типов уроков
// @Summary Получение типов уроков
// @Description Возвращает список всех уникальных типов уроков
// @Tags lessons
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /lessons/types [get]
func (c *LessonController) GetLessonTypes(e echo.Context) error {
	types, err := c.svc.GetLessonTypes()
	if err != nil {
		c.logger.Error("Ошибка получения типов уроков", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.ErrorResponse("SERVER_ERROR", err.Error()))
	}
	return e.JSON(http.StatusOK, responses.DataResponse(map[string]interface{}{
		"types": types,
	}))
}

// GetLessonsByType возвращает список уроков определенного типа
// @Summary Получение уроков по типу
// @Description Возвращает список всех уроков указанного типа
// @Tags lessons
// @Accept json
// @Produce json
// @Param type path string true "Тип урока"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /lessons/type/{type} [get]
func (c *LessonController) GetLessonsByType(e echo.Context) error {
	lessonType := e.Param("type")
	if lessonType == "" {
		return e.JSON(http.StatusBadRequest, responses.ErrorResponse("INVALID_REQUEST", "Тип урока не указан"))
	}

	lessons, err := c.svc.GetLessonsByType(lessonType)
	if err != nil {
		c.logger.Error("Ошибка получения уроков по типу", zap.String("type", lessonType), zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.ErrorResponse("SERVER_ERROR", err.Error()))
	}

	// Преобразование в DTO с парсингом JSON
	lessonDTOs := make([]map[string]interface{}, len(lessons))
	for i, lesson := range lessons {
		var lessonData map[string]interface{}
		if err := json.Unmarshal(lesson.LessonData, &lessonData); err != nil {
			c.logger.Error("Ошибка десериализации данных урока", zap.Error(err))
			continue
		}

		lessonDTOs[i] = map[string]interface{}{
			"id":            lesson.ID,
			"type":          lesson.Type,
			"level":         lesson.Level,
			"mode":          lesson.Mode,
			"english_level": lesson.EnglishLevel,
			"xp":            lesson.XP,
			"createdAt":     lesson.CreatedAt,
			"updatedAt":     lesson.UpdatedAt,
		}
	}

	return e.JSON(http.StatusOK, responses.DataResponse(map[string]interface{}{
		"lessons": lessonDTOs,
	}))
}

// GetLesson возвращает детали конкретного урока
// @Summary Получение детальной информации урока
// @Description Возвращает детальную информацию об уроке по его ID, включая полные данные урока
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "ID урока"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /lessons/{id} [get]
func (c *LessonController) GetLesson(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, responses.ErrorResponse("INVALID_ID", "Некорректный ID урока"))
	}

	lesson, err := c.svc.GetLessonWithParsedData(uint(id))
	if err != nil {
		c.logger.Error("Ошибка получения урока", zap.String("id", idStr), zap.Error(err))
		return e.JSON(http.StatusNotFound, responses.ErrorResponse("LESSON_NOT_FOUND", "Урок не найден"))
	}

	return e.JSON(http.StatusOK, responses.DataResponse(lesson))
}

// GetAllLessons возвращает все уроки
// @Summary Получение всех уроков
// @Description Возвращает список всех уроков
// @Tags lessons
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /lessons [get]
func (c *LessonController) GetAllLessons(e echo.Context) error {
	lessons, err := c.svc.GetAllLessons()
	if err != nil {
		c.logger.Error("Ошибка получения всех уроков", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, responses.ErrorResponse("SERVER_ERROR", err.Error()))
	}

	// Преобразование в DTO с парсингом JSON
	lessonDTOs := make([]map[string]interface{}, len(lessons))
	for i, lesson := range lessons {
		var lessonData map[string]interface{}
		if err := json.Unmarshal(lesson.LessonData, &lessonData); err != nil {
			c.logger.Error("Ошибка десериализации данных урока", zap.Error(err))
			continue
		}

		lessonDTOs[i] = map[string]interface{}{
			"id":            lesson.ID,
			"type":          lesson.Type,
			"level":         lesson.Level,
			"mode":          lesson.Mode,
			"english_level": lesson.EnglishLevel,
			"xp":            lesson.XP,
			"createdAt":     lesson.CreatedAt,
			"updatedAt":     lesson.UpdatedAt,
		}
	}

	return e.JSON(http.StatusOK, responses.DataResponse(map[string]interface{}{
		"lessons": lessonDTOs,
	}))
}
