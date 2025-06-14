package controller

import (
	"net/http"

	"github.com/Ravr-Site/Ravr-Backend/internal/responses"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Типы для Swagger документации (без дженериков)
// @Description Успешный ответ с сообщением
type SwaggerMessageResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    map[string]string `json:"data"`
}

// @Description Успешный ответ с токеном авторизации
type SwaggerTokenResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
}

// @Description Ответ с данными профиля пользователя
type SwaggerUserProfileResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		ID        uint   `json:"id" example:"1"`
		Username  string `json:"username" example:"johndoe"`
		FirstName string `json:"first_name,omitempty" example:"John"`
		LastName  string `json:"last_name,omitempty" example:"Doe"`
	} `json:"data"`
}

// Request structs
// @Description Запрос на регистрацию нового пользователя
type registerRequest struct {
	Username  string `json:"username" validate:"required,min=3" example:"johndoe"`
	Password  string `json:"password" validate:"required,min=6" example:"secret123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
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
// @Success 201 {object} SwaggerMessageResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /auth/register [post]
func (h *UserController) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("VALIDATION_ERROR", err.Error()))
	}
	if err := h.svc.Register(req.Username, req.Password, req.FirstName, req.LastName); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("REGISTRATION_ERROR", err.Error()))
	}
	return c.JSON(http.StatusCreated, responses.MessageResponse("Пользователь успешно зарегистрирован"))
}

// @Description Ответ с токеном авторизации
type TokenResponse struct {
	Token string `json:"token"`
}

// Login аутентифицирует пользователя и выдает JWT токен
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Учетные данные"
// @Success 200 {object} SwaggerTokenResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Router /auth/login [post]
func (h *UserController) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error("VALIDATION_ERROR", err.Error()))
	}
	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.Error("AUTHENTICATION_ERROR", err.Error()))
	}

	return c.JSON(http.StatusOK, responses.Success(TokenResponse{
		Token: token,
	}))
}

// Profile возвращает информацию о текущем пользователе
// @Summary Профиль пользователя
// @Description Возвращает данные текущего аутентифицированного пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerUserProfileResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/user [get]
func (h *UserController) Profile(c echo.Context) error {
	// Получаем username из контекста (установлено JWT middleware)
	username := c.Get("username").(string)

	u, err := h.svc.GetByUsername(username)
	if err != nil {
		h.logger.Error("failed to get user profile",
			zap.String("username", username),
			zap.Error(err))
		return c.JSON(http.StatusInternalServerError, responses.Error("SERVER_ERROR", err.Error()))
	}

	userResp := struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	}{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	return c.JSON(http.StatusOK, responses.Success(userResp))
}
