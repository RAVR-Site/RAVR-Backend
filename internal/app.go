package internal

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/controller"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Application struct {
	container *dig.Container
	config    *config.Config
}

func NewApplication(config *config.Config) *Application {
	app := &Application{
		config:    config,
		container: dig.New(),
	}

	return app
}

func (app *Application) Init() error {
	if err := app.container.Provide(func() (*zap.Logger, error) {
		env := os.Getenv("ENVIRONMENT")
		if env == "" {
			env = "local"
		}

		if env == "local" {
			l, err := zap.NewDevelopment()
			if err != nil {
				return nil, err
			}
			return l, nil
		}

		l, err := zap.NewProduction()
		if err != nil {
			return nil, err // Fallback to no-op logger if production logger fails
		}
		return l, nil
	}); err != nil {
		return err
	}

	if err := app.container.Provide(func(logger *zap.Logger) (*gorm.DB, error) {
		db, err := gorm.Open(postgres.Open(app.config.DatabaseDSN), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		return db, nil
	}); err != nil {
		return err
	}

	if err := app.container.Provide(func() *echo.Echo {
		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(app.config.JWTSecret),
			Skipper: func(c echo.Context) bool {
				// Skip JWT authentication for certain routes
				switch c.Path() {
				case "/api/v1/user/login", "/api/v1/user/register":
					return true
				}
				return false
			},
		}))

		return e
	}); err != nil {
		return err
	}

	if err := app.initRepositories(); err != nil {
		return err
	}

	if err := app.initServices(); err != nil {
		return err
	}

	if err := app.initControllers(); err != nil {
		return err
	}

	return nil
}

func (app *Application) initRepositories() error {
	err := app.container.Provide(repository.NewUserRepository)
	return err
}

func (app *Application) initServices() error {
	err := app.container.Provide(func(repo repository.UserRepository, logger *zap.Logger) service.UserService {
		return service.NewUserService(repo, app.config.JWTSecret, logger)
	})
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) initControllers() error {
	if err := app.container.Invoke(func(e *echo.Echo, userService service.UserService, logger *zap.Logger) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		userHandler := controller.NewUserController(userService, logger)
		e.GET("/api/v1/user", userHandler.Profile)
		e.POST("/api/v1/user/login", userHandler.Login)
		e.POST("/api/v1/user/register", userHandler.Register)
	}); err != nil {
		return err
	}

	return nil
}

func (app *Application) Start() error {
	return app.container.Invoke(func(e *echo.Echo) error {
		return e.Start(":" + app.config.Port)
	})
}
