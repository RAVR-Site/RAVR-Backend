package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ravr-Site/Ravr-Backend/repository"
)

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	GetByUsername(username string) (*repository.User, error)
}

type service struct {
	repo      repository.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

func NewUserService(repo repository.UserRepository, jwtSecret string, logger *zap.Logger) UserService {
	return &service{repo, jwtSecret, logger}
}

func (s *service) Register(username, password string) error {
	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return errors.New("username already taken")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &repository.User{Username: username, Password: string(hash)}
	return s.repo.Create(user)
}

func (s *service) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	signed, err := token.SignedString([]byte(s.jwtSecret))
	s.logger.Info("user logged in", zap.String("username", user.Username))
	return signed, err
}

func (s *service) GetByUsername(username string) (*repository.User, error) {
	return s.repo.FindByUsername(username)
}
