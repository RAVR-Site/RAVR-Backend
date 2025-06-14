package repository

import (
	"time"

	"gorm.io/gorm"
)

// UserRanking хранит информацию о позиции пользователя в рейтинге за определенный период
type UserRanking struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"column:user_id;not null"`
	Position    int       `gorm:"column:position;not null"`
	Experience  uint64    `gorm:"column:experience;not null"`
	Period      string    `gorm:"column:period;not null"` // daily, weekly, monthly
	PeriodStart time.Time `gorm:"column:period_start;not null"`
	PeriodEnd   time.Time `gorm:"column:period_end;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	User        *User     `gorm:"foreignKey:UserID"`
}

// LeaderboardEntry представляет запись в таблице лидеров, включая тренд изменения
type LeaderboardEntry struct {
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Position   int    `json:"position"`
	Experience uint64 `json:"experience"`
	Trend      string `json:"trend"` // "up" - рост, "down" - падение, "stable" - без изменений
}

type LeaderboardRepository interface {
	SaveRankings(rankings []*UserRanking) error
	GetLatestRankings(period string, limit int) ([]*UserRanking, error)
	GetPreviousRankings(period string, beforeDate time.Time, limit int) ([]*UserRanking, error)
	GetUserRankingHistory(userID uint, period string, limit int) ([]*UserRanking, error)
	CalculateLeaderboard(limit int) ([]*LeaderboardEntry, error)
}

type leaderboardRepo struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) LeaderboardRepository {
	return &leaderboardRepo{db}
}

func (r *leaderboardRepo) SaveRankings(rankings []*UserRanking) error {
	return r.db.Create(&rankings).Error
}

func (r *leaderboardRepo) GetLatestRankings(period string, limit int) ([]*UserRanking, error) {
	var rankings []*UserRanking

	// Получаем максимальную дату окончания периода
	var maxPeriodEnd time.Time
	err := r.db.Model(&UserRanking{}).
		Where("period = ?", period).
		Select("MAX(period_end)").
		Row().
		Scan(&maxPeriodEnd)

	if err != nil {
		return nil, err
	}

	if maxPeriodEnd.IsZero() {
		return rankings, nil
	}

	// Получаем рейтинги за последний период
	err = r.db.Where("period = ? AND period_end = ?", period, maxPeriodEnd).
		Order("position ASC").
		Limit(limit).
		Preload("User").
		Find(&rankings).Error

	return rankings, err
}

func (r *leaderboardRepo) GetPreviousRankings(period string, beforeDate time.Time, limit int) ([]*UserRanking, error) {
	var rankings []*UserRanking

	// Находим максимальную дату окончания периода до указанной даты
	var prevPeriodEnd time.Time
	err := r.db.Model(&UserRanking{}).
		Where("period = ? AND period_end < ?", period, beforeDate).
		Select("MAX(period_end)").
		Row().
		Scan(&prevPeriodEnd)

	if err != nil {
		return nil, err
	}

	if prevPeriodEnd.IsZero() {
		return rankings, nil
	}

	// Получаем рейтинги за предыдущий период
	err = r.db.Where("period = ? AND period_end = ?", period, prevPeriodEnd).
		Order("position ASC").
		Limit(limit).
		Preload("User").
		Find(&rankings).Error

	return rankings, err
}

func (r *leaderboardRepo) GetUserRankingHistory(userID uint, period string, limit int) ([]*UserRanking, error) {
	var rankings []*UserRanking

	err := r.db.Where("user_id = ? AND period = ?", userID, period).
		Order("period_end DESC").
		Limit(limit).
		Find(&rankings).Error

	return rankings, err
}

func (r *leaderboardRepo) CalculateLeaderboard(limit int) ([]*LeaderboardEntry, error) {
	// Получаем топ пользователей по опыту
	var users []*User
	err := r.db.Order("experience DESC").Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}

	entries := make([]*LeaderboardEntry, len(users))

	// Получаем рейтинги предыдущего периода, если они есть
	now := time.Now()
	rankings, err := r.GetLatestRankings("weekly", limit*2) // Берем с запасом
	if err != nil {
		return nil, err
	}

	prevRankings, err := r.GetPreviousRankings("weekly", now, limit*2)
	if err != nil {
		return nil, err
	}

	// Создаем карту для быстрого поиска предыдущей позиции
	prevPositions := make(map[uint]int)
	for _, rank := range prevRankings {
		prevPositions[rank.UserID] = rank.Position
	}

	// Заполняем записи таблицы лидеров
	for i, user := range users {
		position := i + 1 // Позиция начинается с 1
		entry := &LeaderboardEntry{
			UserID:     user.ID,
			Username:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Position:   position,
			Experience: user.Experience,
			Trend:      "stable", // По умолчанию без изменений
		}

		// Если есть предыдущая позиция, вычисляем тренд
		if prevPos, exists := prevPositions[user.ID]; exists {
			if prevPos < position {
				entry.Trend = "up" // Рост позиции
			} else if prevPos > position {
				entry.Trend = "down" // Падение позиции
			}
		}

		entries[i] = entry
	}

	return entries, nil
}
