package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockUserService struct{ mock.Mock }

func (m *mockUserService) Register(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}
func (m *mockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
func (m *mockUserService) GetByUsername(username string) (*repository.User, error) {
	args := m.Called(username)
	if u, ok := args.Get(0).(*repository.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("Register", "user", "password").Return(nil)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"user","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		// Проверяем обновлённый формат ответа
		assert.Contains(t, rec.Body.String(), "success")
		assert.Contains(t, rec.Body.String(), "Пользователь успешно зарегистрирован")
	}
}

func TestRegister_ValidationError(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"u","password":"p"}` // слишком коротко
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}

func TestRegister_ServiceError(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("Register", "user", "password").Return(errors.New("username taken"))
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"user","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "username taken")
}

func TestLogin_Success(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("Login", "user", "password").Return("token123", nil)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"user","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "token123")
}

func TestLogin_ValidationError(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}

func TestLogin_ServiceError(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("Login", "user", "password").Return("", errors.New("invalid credentials"))
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	body := `{"username":"user","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid credentials")
}

func TestProfile_Success(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("GetByUsername", "user").Return(&repository.User{ID: 1, Username: "user"}, nil)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	// Создаем правильный HTTP запрос и рекордер
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Устанавливаем данные пользователя в контекст, как это делает JWT middleware
	c.Set("username", "user")
	c.Set("user_id", 1)

	// Вызываем тестируемый метод
	err := h.Profile(c)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "user")    // Имя пользователя должно быть в ответе
	assert.Contains(t, rec.Body.String(), "success") // Новый формат ответа
}

func TestProfile_Error(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("GetByUsername", "user").Return(nil, errors.New("db error"))
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	// Создаем правильный HTTP запрос и рекордер
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Устанавливаем данные пользователя в контекст, как это делает JWT middleware
	c.Set("username", "user")
	c.Set("user_id", 1)

	// Вызываем тестируемый метод
	err := h.Profile(c)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "error") // Проверяем, что ответ содержит информацию об ошибке
}

func TestUploadImage_Success(t *testing.T) {
	t.Skip("Upload image functionality needs proper storage implementation")
}

func TestUploadImage_SaveError(t *testing.T) {
	t.Skip("Upload image functionality needs proper storage implementation")
}
