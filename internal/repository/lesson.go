package repository

import (
	"gorm.io/gorm"
	"time"
)

// Lesson представляет урок в базе данных
type Lesson struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Type         string    `gorm:"not null" json:"type"`
	Level        string    `gorm:"not null" json:"level"`
	Mode         string    `gorm:"not null" json:"mode"`
	EnglishLevel string    `gorm:"not null" json:"english_level"` // Уровень владения английским (A1, A2, B1, B2, C1, C2)
	XP           int       `gorm:"not null" json:"xp"`
	LessonData   []byte    `gorm:"type:jsonb;not null" json:"lesson_data"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// LessonRepository интерфейс для работы с уроками
type LessonRepository interface {
	Create(lesson *Lesson) error
	GetByID(id uint) (*Lesson, error)
	GetAll() ([]*Lesson, error)
	GetAllByType(lessonType string) ([]*Lesson, error)
	GetUniqueTypes() ([]string, error)
	Update(lesson *Lesson) error
	Delete(id uint) error
}

// lessonRepo имплементация LessonRepository
type lessonRepo struct {
	db *gorm.DB
}

// NewLessonRepository создает новый экземпляр LessonRepository
func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepo{db}
}

// Create создает новый урок в базе данных
func (r *lessonRepo) Create(lesson *Lesson) error {
	return r.db.Create(lesson).Error
}

// GetByID возвращает урок по его ID
func (r *lessonRepo) GetByID(id uint) (*Lesson, error) {
	var lesson Lesson
	err := r.db.First(&lesson, id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// GetAll возвращает все уроки
func (r *lessonRepo) GetAll() ([]*Lesson, error) {
	var lessons []*Lesson
	err := r.db.Find(&lessons).Error
	return lessons, err
}

// GetAllByType возвращает все уроки определенного типа
func (r *lessonRepo) GetAllByType(lessonType string) ([]*Lesson, error) {
	var lessons []*Lesson
	err := r.db.Where("type = ?", lessonType).Find(&lessons).Error
	return lessons, err
}

// GetUniqueTypes возвращает список уникальных типов уроков
func (r *lessonRepo) GetUniqueTypes() ([]string, error) {
	var types []string
	err := r.db.Model(&Lesson{}).Distinct().Pluck("type", &types).Error
	return types, err
}

// Update обновляет информацию об уроке
func (r *lessonRepo) Update(lesson *Lesson) error {
	return r.db.Save(lesson).Error
}

// Delete удаляет урок по ID
func (r *lessonRepo) Delete(id uint) error {
	return r.db.Delete(&Lesson{}, id).Error
}
