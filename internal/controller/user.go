package controller

import (
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/Ravr-Site/Ravr-Backend/internal/storage"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Request structs
type registerRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserController struct {
	svc      service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserController(svc service.UserService, logger *zap.Logger) *UserController {
	return &UserController{svc, validator.New(), logger}
}

// Register @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true
// @Success 201 {object} map[string]string
// @Router /register [post]
func (h *UserController) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := h.svc.Register(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "registered"})
}

// Login @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true
// @Success 200 {object} map[string]string
// @Router /login [post]
func (h *UserController) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *UserController) Profile(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	u, err := h.svc.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"id": u.ID, "username": u.Username})
}

func (h *UserController) UploadImage(store storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid file"})
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		filename := filepath.Base(file.Filename)
		dstPath := filepath.Join("uploads", filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()
		io.Copy(dst, src)

		// store via storage interface
		if err := store.Save("uploads", filename); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"filename": filename})
	}
}

func ServeImage(c echo.Context) error {
	filename := c.Param("filename")
	return c.File(filepath.Join("uploads", filename))
}
