package service

import (
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
)

type ResultService interface {
	Save(username string, lessonId, timeTaken uint) error
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

func (s *resultService) Save(username string, lessonId, timeTaken uint) error {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return err
	}
	result := &repository.Result{
		UserID:    user.ID,
		LessonID:  lessonId,
		TimeTaken: timeTaken,
	}

	if err = s.repo.Create(result); err != nil {
		s.logger.Error("Failed to save result", zap.Error(err))
		return err
	}

	s.logger.Info("Result saved successfully", zap.Uint("user_id", user.ID), zap.Uint("lesson_id", lessonId), zap.Uint("time_taken", timeTaken))

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

	// Получаем результаты вокруг пользователя из репозитория
	results, userPosition, err := s.repo.GetLeaderboardAroundUser(user.ID, lessonId, limit)
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

		leaderboard.Results[i] = Result{
			Position:  position,
			Username:  resultUser.Username,
			TimeTaken: result.TimeTaken,
			XP:        result.XP,
		}
	}

	return leaderboard, nil
}
