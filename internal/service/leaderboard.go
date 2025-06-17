package service

import (
	"time"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
)

// LeaderboardService интерфейс для работы с рейтингом пользователей
type LeaderboardService interface {
	GetLeaderboard(limit int) ([]*repository.LeaderboardEntry, error)
	GetExtendedLeaderboard(limit int) ([]*repository.ExtendedLeaderboardEntry, error)
	GetLessonLeaderboard(lessonID string, userID uint, limit int) ([]*repository.LessonLeaderboardEntry, int, error)
	UpdateUserRankings(period string) error
}

type leaderboardService struct {
	userRepo        repository.UserRepository
	leaderboardRepo repository.LeaderboardRepository
	logger          *zap.Logger
}

// NewLeaderboardService создает новый сервис для работы с рейтингом пользователей
func NewLeaderboardService(
	userRepo repository.UserRepository,
	leaderboardRepo repository.LeaderboardRepository,
	logger *zap.Logger,
) LeaderboardService {
	return &leaderboardService{
		userRepo:        userRepo,
		leaderboardRepo: leaderboardRepo,
		logger:          logger,
	}
}

// GetLeaderboard возвращает текущую таблицу лидеров с трендами изменения позиций
func (s *leaderboardService) GetLeaderboard(limit int) ([]*repository.LeaderboardEntry, error) {
	entries, err := s.leaderboardRepo.CalculateLeaderboard(limit)
	if err != nil {
		s.logger.Error("Failed to calculate leaderboard", zap.Error(err))
		return nil, err
	}
	return entries, nil
}

// GetExtendedLeaderboard возвращает расширенную таблицу лидеров с информацией о времени
func (s *leaderboardService) GetExtendedLeaderboard(limit int) ([]*repository.ExtendedLeaderboardEntry, error) {
	entries, err := s.leaderboardRepo.CalculateExtendedLeaderboard(limit)
	if err != nil {
		s.logger.Error("Failed to calculate extended leaderboard", zap.Error(err))
		return nil, err
	}
	return entries, nil
}

// GetLessonLeaderboard возвращает таблицу лидеров для конкретного урока и позицию пользователя в ней
func (s *leaderboardService) GetLessonLeaderboard(lessonID string, userID uint, limit int) ([]*repository.LessonLeaderboardEntry, int, error) {
	entries, userPosition, err := s.leaderboardRepo.GetLessonLeaderboard(lessonID, userID, limit)
	if err != nil {
		s.logger.Error("Failed to get lesson leaderboard",
			zap.Error(err),
			zap.String("lessonID", lessonID),
			zap.Uint("userID", userID))
		return nil, 0, err
	}
	return entries, userPosition, nil
}

// UpdateUserRankings обновляет рейтинги пользователей за указанный период (weekly, monthly)
func (s *leaderboardService) UpdateUserRankings(period string) error {
	// Получаем топ пользователей по опыту
	users, err := s.userRepo.GetTopUsersByExperience(1000) // Берем с запасом
	if err != nil {
		s.logger.Error("Failed to get top users", zap.Error(err))
		return err
	}

	// Определяем даты начала и конца периода
	now := time.Now()
	var periodStart, periodEnd time.Time

	switch period {
	case "daily":
		periodStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 0, 1)
	case "weekly":
		// Находим начало недели (понедельник)
		daysFromMonday := int(now.Weekday())
		if daysFromMonday == 0 { // Воскресенье
			daysFromMonday = 7
		}
		daysFromMonday--
		periodStart = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 0, 7)
	case "monthly":
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	default:
		s.logger.Error("Invalid period specified", zap.String("period", period))
		return nil
	}

	// Создаем записи рейтинга
	rankings := make([]*repository.UserRanking, len(users))
	for i, user := range users {
		rankings[i] = &repository.UserRanking{
			UserID:      user.ID,
			Position:    i + 1, // Позиции начинаются с 1
			Experience:  user.Experience,
			Period:      period,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}
	}

	// Сохраняем рейтинги
	err = s.leaderboardRepo.SaveRankings(rankings)
	if err != nil {
		s.logger.Error("Failed to save rankings", zap.Error(err))
		return err
	}

	s.logger.Info("User rankings updated successfully",
		zap.String("period", period),
		zap.Int("users_count", len(users)),
		zap.Time("period_start", periodStart),
		zap.Time("period_end", periodEnd))

	return nil
}
