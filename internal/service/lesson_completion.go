package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
)

// LessonCompletionResult содержит результат завершения урока
type LessonCompletionResult struct {
	UserID           uint
	LessonID         string
	CompletionTime   uint64 // Время в секундах
	EarnedExperience uint64
	TotalExperience  uint64
}

// LessonCompletionService интерфейс для работы с завершением уроков
type LessonCompletionService interface {
	CompleteLessonAndUpdateExperience(userID uint, lessonID string, completionTime uint64, earnedExperience uint64) (*LessonCompletionResult, error)
}

type lessonCompletionService struct {
	userRepo   repository.UserRepository
	resultRepo repository.ResultRepository
	lessonRepo repository.LessonRepository
	logger     *zap.Logger
}

// NewLessonCompletionService создает новый сервис для работы с завершением уроков
func NewLessonCompletionService(
	userRepo repository.UserRepository,
	resultRepo repository.ResultRepository,
	lessonRepo repository.LessonRepository,
	logger *zap.Logger,
) LessonCompletionService {
	return &lessonCompletionService{
		userRepo:   userRepo,
		resultRepo: resultRepo,
		lessonRepo: lessonRepo,
		logger:     logger,
	}
}

// CompleteLessonAndUpdateExperience обрабатывает завершение урока,
// обновляет опыт пользователя и записывает в историю результатов
func (s *lessonCompletionService) CompleteLessonAndUpdateExperience(
	userID uint,
	lessonID string,
	completionTime uint64,
	earnedExperience uint64,
) (*LessonCompletionResult, error) {
	// Получаем пользователя
	userList, err := s.userRepo.GetUsersByIds([]uint{userID})
	if err != nil {
		s.logger.Error("Error getting user", zap.Error(err), zap.Uint("userId", userID))
		return nil, err
	}

	if len(userList) == 0 {
		return nil, errors.New("user not found")
	}

	user := userList[0]

	// Конвертируем строковый ID урока в uint
	lessonIDUint, err := strconv.ParseUint(lessonID, 10, 32)
	if err != nil {
		s.logger.Error("Error converting lesson ID to uint", zap.Error(err), zap.String("lessonId", lessonID))
		return nil, fmt.Errorf("invalid lesson ID format: %v", err)
	}

	// Получаем урок для проверки его существования
	lesson, err := s.lessonRepo.GetByID(uint(lessonIDUint))
	if err != nil {
		s.logger.Error("Error getting lesson", zap.Error(err), zap.String("lessonId", lessonID))
		return nil, err
	}

	if lesson == nil {
		return nil, errors.New("lesson not found")
	}

	// Обновляем опыт пользователя
	newExperience := user.Experience + earnedExperience
	err = s.userRepo.UpdateExperience(userID, newExperience)
	if err != nil {
		s.logger.Error("Error updating user experience", zap.Error(err), zap.Uint("userId", userID))
		return nil, err
	}

	now := time.Now()

	// Проверяем, существует ли уже результат для данного пользователя и урока
	existingResult, err := s.resultRepo.GetResultByUserAndLesson(userID, lessonID)

	// Если запись существует, обновляем её
	if err == nil && existingResult != nil {
		existingResult.Score = earnedExperience // Обновляем счет
		existingResult.CompletedAt = now        // Обновляем время завершения
		existingResult.CompletionTime = completionTime
		existingResult.AddedExperience = earnedExperience // Обновляем добавленный опыт
		existingResult.XP = earnedExperience              // Обновляем XP для обратной совместимости

		err = s.resultRepo.Update(existingResult)
		if err != nil {
			s.logger.Error("Error updating existing result", zap.Error(err), zap.Uint("userId", userID), zap.String("lessonId", lessonID))
			return nil, err
		}

		s.logger.Info("Existing lesson result updated",
			zap.Uint("userId", userID),
			zap.String("lessonId", lessonID),
			zap.Uint64("completionTime", completionTime),
			zap.Uint64("earnedExperience", earnedExperience))
	} else {
		// Если записи нет, создаем новую
		result := &repository.Result{
			UserID:          userID,
			LessonID:        lessonID,
			Score:           earnedExperience, // Используем заработанный опыт как счет
			CompletedAt:     now,
			CompletionTime:  completionTime, // Сохраняем время завершения в секундах
			AddedExperience: earnedExperience,
			XP:              earnedExperience, // Заполняем устаревшее поле для обратной совместимости
		}

		err = s.resultRepo.Create(result)
		if err != nil {
			s.logger.Error("Error creating result record", zap.Error(err), zap.Uint("userId", userID), zap.String("lessonId", lessonID))
			return nil, err
		}

		s.logger.Info("New lesson result created",
			zap.Uint("userId", userID),
			zap.String("lessonId", lessonID),
			zap.Uint64("completionTime", completionTime),
			zap.Uint64("earnedExperience", earnedExperience))
	}

	// Формируем результат
	completionResult := &LessonCompletionResult{
		UserID:           userID,
		LessonID:         lessonID,
		CompletionTime:   completionTime,
		EarnedExperience: earnedExperience,
		TotalExperience:  newExperience,
	}

	s.logger.Info("Lesson completed successfully",
		zap.Uint("userId", userID),
		zap.String("lessonId", lessonID),
		zap.Uint64("completionTime", completionTime),
		zap.Uint64("earnedExperience", earnedExperience),
		zap.Uint64("totalExperience", newExperience))

	return completionResult, nil
}
