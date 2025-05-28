package controller

import (
	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/Ravr-Site/Ravr-Backend/internal/storage"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Request structs
// @Description Запрос на регистрацию нового пользователя
type registerRequest struct {
	Username string `json:"username" validate:"required,min=3" example:"johndoe"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// @Description Запрос на вход в систему
type loginRequest struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

type UserController struct {
	svc      service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserController(svc service.UserService, logger *zap.Logger) *UserController {
	return &UserController{svc, validator.New(), logger}
}

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя с указанными учетными данными
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Данные для регистрации"
// @Success 201 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Router /user/register [post]
func (h *UserController) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("INVALID_REQUEST", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("VALIDATION_ERROR", err.Error()))
	}
	if err := h.svc.Register(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("REGISTRATION_ERROR", err.Error()))
	}
	return c.JSON(http.StatusCreated, responses.MessageResponse("Пользователь успешно зарегистрирован"))
}

// Login аутентифицирует пользователя и выдает JWT токен
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Учетные данные"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Router /user/login [post]
func (h *UserController) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("INVALID_REQUEST", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("VALIDATION_ERROR", err.Error()))
	}
	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse("AUTHENTICATION_ERROR", err.Error()))
	}
	return c.JSON(http.StatusOK, responses.DataResponse(responses.TokenResponse{Token: token}))
}

// Profile возвращает информацию о текущем пользователе
// @Summary Профиль пользователя
// @Description Возвращает данные текущего аутентифицированного пользователя
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /user [get]
func (h *UserController) Profile(c echo.Context) error {
	// Получаем username из контекста (установлено JWT middleware)
	username := c.Get("username").(string)

	u, err := h.svc.GetByUsername(username)
	if err != nil {
		h.logger.Error("failed to get user profile",
			zap.String("username", username),
			zap.Error(err))
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse("SERVER_ERROR", err.Error()))
	}

	return c.JSON(http.StatusOK, responses.DataResponse(responses.UserResponse{
		ID:       u.ID,
		Username: u.Username,
	}))
}

func (h *UserController) UploadImage(store storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse("INVALID_FILE", "Некорректный файл"))
		}
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse("FILE_UPLOAD_ERROR", err.Error()))
		}
		defer src.Close()

		filename := filepath.Base(file.Filename)
		dstPath := filepath.Join("uploads", filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse("FILE_CREATION_ERROR", err.Error()))
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse("FILE_COPY_ERROR", err.Error()))
		}

		// store via storage interface
		if err := store.Save("uploads", filename); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse("STORAGE_ERROR", err.Error()))
		}
		return c.JSON(http.StatusOK, responses.DataResponse(map[string]string{"filename": filename}))
	}
}

func ServeImage(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(filepath.Join("uploads", filename))
}
