package repository

import (
	"time"

	"gorm.io/gorm"
)

// UserStats содержит статистику пользователя
type UserStats struct {
	TotalLessons      int     // Общее количество пройденных уроков
	TotalExperience   uint64  // Общий полученный опыт
	AverageExperience float64 // Средний опыт за урок
	MaxExperience     uint64  // Максимальный опыт за урок
	FastestCompletion string  // Самое быстрое время завершения (формат MM:SS)
	AverageCompletion float64 // Среднее время завершения (в секундах)
}

type ResultRepository interface {
	Create(result *Result) error
	GetLeaderboardAroundUser(userID uint, lessonID string, limit uint) ([]Result, int, error)
	GetUserStats(userID uint) (*UserStats, error)
}

type resultRepo struct {
	db *gorm.DB
}

func NewResultRepository(db *gorm.DB) ResultRepository {
	return &resultRepo{db}
}

type Result struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"column:user_id;not null"`
	LessonID        string    `gorm:"column:lesson_id;not null"`
	Score           uint64    `gorm:"column:score;not null"`
	CompletedAt     time.Time `gorm:"column:completed_at;not null"`
	CompletionTime  string    `gorm:"column:completion_time"` // Формат MM:SS или HH:MM:SS
	AddedExperience uint64    `gorm:"column:added_experience"`
}

func (r *resultRepo) Create(result *Result) error {
	return r.db.Create(result).Error
}

// GetLeaderboardAroundUser получает результаты вокруг заданного пользователя для определенного урока
// Возвращает список результатов, позицию пользователя в списке и ошибку, если она возникла
func (r *resultRepo) GetLeaderboardAroundUser(userID uint, lessonID string, limit uint) ([]Result, int, error) {
	// Сначала найдем позицию пользователя в общем рейтинге для данного урока
	var userPosition int64
	subQuery := r.db.Model(&Result{}).
		Where("lesson_id = ? AND score < (SELECT score FROM results WHERE user_id = ? AND lesson_id = ? LIMIT 1)",
			lessonID, userID, lessonID).
		Count(&userPosition)

	if subQuery.Error != nil {
		return nil, 0, subQuery.Error
	}
	userPosition++ // Позиция начинается с 1, а не с 0

	// Теперь получим результаты вокруг пользователя
	var results []Result

	// Вычисляем начальную позицию для запроса
	startPosition := 1
	if userPosition > int64(limit) {
		startPosition = int(uint(userPosition) - limit)
	}

	// Получаем результаты
	err := r.db.Model(&Result{}).
		Where("lesson_id = ?", lessonID).
		Order("score ASC").
		Limit(int(limit*2 + 1)). // limit выше + limit ниже + результат пользователя
		Offset(startPosition - 1).
		Find(&results).Error

	return results, int(userPosition), err
}

// GetUserStats возвращает статистику пользователя на основе его результатов
func (r *resultRepo) GetUserStats(userID uint) (*UserStats, error) {
	var stats UserStats

	// Получаем общее количество пройденных уроков
	err := r.db.Model(&Result{}).Where("user_id = ?", userID).Count(&stats.TotalLessons).Error
	if err != nil {
		return nil, err
	}

	// Если уроков нет, возвращаем пустую статистику
	if stats.TotalLessons == 0 {
		return &stats, nil
	}

	// Получаем общий опыт
	err = r.db.Model(&Result{}).
		Select("COALESCE(SUM(added_experience), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.TotalExperience).Error
	if err != nil {
		return nil, err
	}

	// Получаем максимальный опыт за урок
	err = r.db.Model(&Result{}).
		Select("COALESCE(MAX(added_experience), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.MaxExperience).Error
	if err != nil {
		return nil, err
	}

	// Вычисляем средний опыт
	if stats.TotalLessons > 0 {
		stats.AverageExperience = float64(stats.TotalExperience) / float64(stats.TotalLessons)
	}

	// Получаем самое быстрое время завершения
	var fastestResult Result
	err = r.db.Where("user_id = ?", userID).
		Order("score").Limit(1).Find(&fastestResult).Error
	if err != nil {
		return nil, err
	}
	stats.FastestCompletion = fastestResult.CompletionTime

	// Получаем среднее время завершения (в секундах)
	err = r.db.Model(&Result{}).
		Select("COALESCE(AVG(score), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.AverageCompletion).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
