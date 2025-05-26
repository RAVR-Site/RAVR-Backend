package controller

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/golang-jwt/jwt/v4"
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

type mockStorage struct{ mock.Mock }

func (m *mockStorage) Save(folder, filename string) error {
	args := m.Called(folder, filename)
	return args.Error(0)
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
		assert.Contains(t, rec.Body.String(), "registered")
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": "user"})
	c := e.NewContext(nil, httptest.NewRecorder())
	c.Set("user", token)

	err := h.Profile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, c.Response().Status)
}

func TestProfile_Error(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	ms.On("GetByUsername", "user").Return(nil, errors.New("db error"))
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": "user"})
	c := e.NewContext(nil, httptest.NewRecorder())
	c.Set("user", token)

	err := h.Profile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, c.Response().Status)
}

func TestUploadImage_Success(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	store := new(mockStorage)
	store.On("Save", "uploads", "test.png").Return(nil)
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	file, _ := os.CreateTemp("", "test-*.png")
	defer os.Remove(file.Name())
	file.WriteString("data")
	file.Seek(0, 0)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("image", "test.png")
	io.Copy(fw, file)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.UploadImage(store)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test.png")
}

func TestUploadImage_SaveError(t *testing.T) {
	e := echo.New()
	ms := new(mockUserService)
	store := new(mockStorage)
	store.On("Save", "uploads", "test.png").Return(errors.New("save error"))
	logger := zap.NewNop()
	h := NewUserController(ms, logger)

	file, _ := os.CreateTemp("", "test-*.png")
	defer os.Remove(file.Name())
	file.WriteString("data")
	file.Seek(0, 0)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("image", "test.png")
	io.Copy(fw, file)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.UploadImage(store)(c)
	assert.Error(t, err)
}
