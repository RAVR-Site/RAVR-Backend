package service

import (
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(user *repository.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *mockUserRepo) FindByUsername(username string) (*repository.User, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*repository.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserService_Register_Success(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByUsername", "user").Return(nil, nil)
	repo.On("Create", mock.AnythingOfType("*repository.User")).Return(nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	err := svc.Register("user", "pass")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_Register_UsernameTaken(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByUsername", "user").Return(&repository.User{Username: "user"}, nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	err := svc.Register("user", "pass")
	assert.Error(t, err)
	assert.Equal(t, "username already taken", err.Error())
}

func TestUserService_Login_Success(t *testing.T) {
	repo := new(mockUserRepo)
	user := &repository.User{Username: "user"}
	_ = svcPasswordHash(user, "pass")
	repo.On("FindByUsername", "user").Return(user, nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	token, err := svc.Login("user", "pass")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	parsed, _ := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) { return []byte("secret"), nil })
	assert.True(t, parsed.Valid)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	repo := new(mockUserRepo)
	user := &repository.User{Username: "user"}
	_ = svcPasswordHash(user, "pass")
	repo.On("FindByUsername", "user").Return(user, nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	_, err := svc.Login("user", "wrong")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByUsername", "nouser").Return(nil, nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	_, err := svc.Login("nouser", "pass")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestUserService_GetByUsername(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByUsername", "user").Return(&repository.User{Username: "user"}, nil)
	logger := zap.NewNop()
	svc := NewUserService(repo, "secret", 24, logger)

	u, err := svc.GetByUsername("user")
	assert.NoError(t, err)
	assert.Equal(t, "user", u.Username)
}

// helper for password hash in tests
func svcPasswordHash(u *repository.User, pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
