package service

import (
	"errors"

	"github.com/Ravr-Site/Ravr-Backend/internal/auth"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, password, firstName, lastName string) error
	Login(username, password string) (string, error)
	GetByUsername(username string) (*repository.User, error)
	UpdateUser(username string, firstName, lastName string) error
}

type service struct {
	repo       repository.UserRepository
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

func NewUserService(repo repository.UserRepository, jwtSecret string, jwtAccessExpiration int, logger *zap.Logger) UserService {
	jwtManager := auth.NewJWTManager(jwtSecret, jwtAccessExpiration)
	return &service{repo, jwtManager, logger}
}

func (s *service) Register(username, password, firstName, lastName string) error {
	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return errors.New("username already taken")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &repository.User{
		Username:  username,
		Password:  string(hash),
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.repo.Create(user)
}

func (s *service) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("неверный пароль")
	}

	// Генерируем токен с помощью JWT менеджера
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("failed to generate JWT token", zap.Error(err))
		return "", errors.New("failed to generate token")
	}

	s.logger.Info("user logged in",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID))

	return token, nil
}

func (s *service) GetByUsername(username string) (*repository.User, error) {
	return s.repo.FindByUsername(username)
}

// UpdateUser обновляет данные пользователя
func (s *service) UpdateUser(username string, firstName, lastName string) error {
	// Проверяем, существует ли пользователь
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return errors.New("пользователь не найден")
	}

	// Подготавливаем данные для обновления
	userData := map[string]interface{}{}

	if firstName != "" {
		userData["first_name"] = firstName
	}

	if lastName != "" {
		userData["last_name"] = lastName
	}

	// Если нет данных для обновления, возвращаем nil
	if len(userData) == 0 {
		return nil
	}

	// Обновляем пользователя в репозитории
	err = s.repo.Update(username, userData)
	if err != nil {
		s.logger.Error("failed to update user",
			zap.String("username", username),
			zap.Error(err))
		return errors.New("ошибка при обновлении данных пользователя")
	}

	return nil
}
