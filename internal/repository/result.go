package repository

import (
	"time"

	"gorm.io/gorm"
)

type ResultRepository interface {
	Create(result *Result) error
	GetLeaderboardAroundUser(userID, lessonID, limit uint) ([]Result, int, error)
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
func (r *resultRepo) GetLeaderboardAroundUser(userID, lessonID, limit uint) ([]Result, int, error) {
	// Сначала найдем позицию пользователя в общем рейтинге для данного урока
	var userPosition int64
	subQuery := r.db.Model(&Result{}).
		Where("lesson_id = ? AND time_taken < (SELECT time_taken FROM results WHERE user_id = ? AND lesson_id = ? LIMIT 1)",
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
		Order("time_taken ASC").
		Limit(int(limit*2 + 1)). // limit выше + limit ниже + результат пользователя
		Offset(startPosition - 1).
		Find(&results).Error

	return results, int(userPosition), err
}
