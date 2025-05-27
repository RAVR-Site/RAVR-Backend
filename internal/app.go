package internal

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/controller"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
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

		// Устанавливаем обработчик 404 ошибок
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			code := http.StatusInternalServerError
			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
			}

			// Если маршрут не найден, возвращаем 404
			if code == http.StatusNotFound {
				_ = c.JSON(http.StatusNotFound, map[string]string{
					"message": "Endpoint not found",
				})
				return
			}

			// Для остальных ошибок используем стандартный обработчик
			e.DefaultHTTPErrorHandler(err, c)
		}

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
		// Маршруты, доступные без JWT аутентификации
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.GET("swagger/doc.json", func(c echo.Context) error {
			return c.File("docs/doc.json")
		})
		e.GET("swagger/doc.yaml", func(c echo.Context) error {
			return c.File("docs/doc.yaml")
		})

		userHandler := controller.NewUserController(userService, logger)
		e.POST("/api/v1/user/login", userHandler.Login)
		e.POST("/api/v1/user/register", userHandler.Register)

		// JWT middleware для защищенных маршрутов
		jwtConfig := echojwt.Config{
			SigningKey: []byte(app.config.JWTSecret),
		}

		// Создаем группу для защищенных маршрутов
		api := e.Group("/api/v1")
		api.Use(echojwt.WithConfig(jwtConfig))

		// Защищенные маршруты
		api.GET("/user", userHandler.Profile)

		// Тут можно добавлять другие защищенные маршруты
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
