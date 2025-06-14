package service

import (
	"fmt"
	"time"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
)

type ResultService interface {
	Save(username string, lessonId string, timeTaken uint) error
	GetLeaderboard(username string, lessonId uint, limit uint) (*Leaderboard, error)
}

type Result struct {
	Position  int
	Username  string
	TimeTaken uint
	XP        uint
}

type Leaderboard struct {
	UserPosition int
	Results      []Result
}

type resultService struct {
	repo        repository.ResultRepository
	userRepo    repository.UserRepository
	userService UserService
	logger      *zap.Logger
}

func NewResultService(repo repository.ResultRepository, userRepo repository.UserRepository, userService UserService, logger *zap.Logger) ResultService {
	return &resultService{
		repo:        repo,
		userRepo:    userRepo,
		userService: userService,
		logger:      logger,
	}
}

func (s *resultService) Save(username string, lessonId string, timeTaken uint) error {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return err
	}

	// Преобразуем время в строку в формате MM:SS
	minutes := timeTaken / 60
	seconds := timeTaken % 60
	completionTime := fmt.Sprintf("%02d:%02d", minutes, seconds)

	result := &repository.Result{
		UserID:         user.ID,
		LessonID:       lessonId,
		Score:          uint64(timeTaken), // Используем timeTaken как Score
		CompletionTime: completionTime,    // Время в формате строки
		CompletedAt:    time.Now(),        // Текущее время
	}

	if err = s.repo.Create(result); err != nil {
		s.logger.Error("Failed to save result", zap.Error(err))
		return err
	}

	s.logger.Info("Result saved successfully", zap.Uint("user_id", user.ID), zap.String("lesson_id", lessonId), zap.Uint("time_taken", timeTaken))

	return nil
}

// GetLeaderboard возвращает таблицу лидеров вокруг указанного пользователя
// limit определяет количество результатов выше и ниже пользователя
func (s *resultService) GetLeaderboard(username string, lessonId uint, limit uint) (*Leaderboard, error) {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	// Конвертируем uint lessonId в string для совместимости с новой структурой Result
	lessonIdStr := fmt.Sprintf("%d", lessonId)

	// Получаем результаты вокруг пользователя из репозитория
	// Обратите внимание: передаем lessonIdStr вместо lessonId, т.к. тип изменился
	results, userPosition, err := s.repo.GetLeaderboardAroundUser(user.ID, lessonIdStr, limit)
	if err != nil {
		s.logger.Error("Failed to get leaderboard", zap.Error(err))
		return nil, err
	}

	// Создаем и заполняем лидерборд
	leaderboard := &Leaderboard{
		UserPosition: userPosition,
		Results:      make([]Result, len(results)),
	}

	// Собираем все ID пользователей для получения информации о них одним запросом
	userIDs := make([]uint, 0, len(results))
	for _, result := range results {
		userIDs = append(userIDs, result.UserID)
	}

	// Получаем всех пользователей одним запросом
	users, err := s.userRepo.GetUsersByIds(userIDs)
	if err != nil {
		s.logger.Error("Failed to get users for leaderboard", zap.Error(err))
		return nil, err
	}

	// Создаем кэш пользователей для быстрого доступа
	userCache := make(map[uint]*repository.User, len(users))
	for _, u := range users {
		userCache[u.ID] = u
	}

	// Заполняем результаты в лидерборде
	for i, result := range results {
		// Определяем позицию в рейтинге
		position := userPosition - (len(results) / 2) + i
		if position < 1 {
			position = i + 1 // Если мы в начале списка, считаем позицию по индексу
		}

		// Получаем пользователя из кэша
		resultUser, exists := userCache[result.UserID]
		if !exists {
			s.logger.Error("User not found in cache", zap.Uint("user_id", result.UserID))
			continue
		}

		// Извлекаем время из поля Score
		timeTaken := uint(result.Score)

		// Используем значение Score как опыт (XP)
		xp := uint(result.Score)

		leaderboard.Results[i] = Result{
			Position:  position,
			Username:  resultUser.Username,
			TimeTaken: timeTaken,
			XP:        xp,
		}
	}

	return leaderboard, nil
}
