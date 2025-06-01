package repository

import (
	"gorm.io/gorm"
	"time"
)

// Lesson представляет урок в базе данных
type Lesson struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `gorm:"not null" json:"type"`
	Level      uint      `gorm:"not null" json:"level"`
	Mode       string    `gorm:"not null" json:"mode"`
	LessonData []byte    `gorm:"type:jsonb;not null" json:"lesson_data"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// LessonRepository интерфейс для работы с уроками
type LessonRepository interface {
	GetAll() ([]*Lesson, error)
	Create(lesson *Lesson) error
	GetByID(id uint) (*Lesson, error)
	GetByType(lessonType string) ([]*Lesson, error)
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

// GetAll возвращает все уроки из базы данных
func (r *lessonRepo) GetAll() ([]*Lesson, error) {
	var lessons []*Lesson
	err := r.db.Find(&lessons).Error
	if err != nil {
		return nil, err
	}
	return lessons, nil
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

// GetByTypeMode возвращает уроки определенного типа
func (r *lessonRepo) GetByType(lessonType string) ([]*Lesson, error) {
	var lessons []*Lesson
	err := r.db.Where("type = ?", lessonType).Find(&lessons).Error
	return lessons, err
}

// Update обновляет информацию об уроке
func (r *lessonRepo) Update(lesson *Lesson) error {
	return r.db.Save(lesson).Error
}

// Delete удаляет урок по ID
func (r *lessonRepo) Delete(id uint) error {
	return r.db.Delete(&Lesson{}, id).Error
}
