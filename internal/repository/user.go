package repository

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	Experience uint64 `gorm:"column:experience;default:0"` // Опыт пользователя
	AvatarURL  string `gorm:"column:avatar_url"`           // URL аватарки пользователя
}

type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	GetUsersByIds(userIDs []uint) ([]*User, error)
	GetTopUsersByExperience(limit int) ([]*User, error)
	UpdateExperience(userID uint, experience uint64) error
	Update(username string, userData map[string]interface{}) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetUsersByIds получает список пользователей по их идентификаторам одним запросом
// Это оптимизированный метод, который избегает повторных запросов к базе данных
func (r *userRepo) GetUsersByIds(userIDs []uint) ([]*User, error) {
	var users []*User

	// Если список ID пустой, просто возвращаем пустой массив
	if len(userIDs) == 0 {
		return users, nil
	}

	// Получаем всех пользователей с указанными ID одним запросом
	err := r.db.Where("id IN ?", userIDs).Find(&users).Error
	return users, err
}

func (r *userRepo) GetTopUsersByExperience(limit int) ([]*User, error) {
	var users []*User
	err := r.db.Order("experience desc").Limit(limit).Find(&users).Error
	return users, err
}

func (r *userRepo) UpdateExperience(userID uint, experience uint64) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("experience", experience).Error
}

// Update обновляет данные пользователя по имени пользователя
func (r *userRepo) Update(username string, userData map[string]interface{}) error {
	return r.db.Model(&User{}).Where("username = ?", username).Updates(userData).Error
}
