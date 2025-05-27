package internal

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/controller"
	"github.com/Ravr-Site/Ravr-Backend/internal/middleware"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"github.com/Ravr-Site/Ravr-Backend/internal/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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

	if err := app.container.Provide(func(_ *zap.Logger) (*gorm.DB, error) {
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

		e.Use(echomiddleware.Logger())
		e.Use(echomiddleware.Recover())

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
	if err := app.container.Provide(repository.NewUserRepository); err != nil {
		return err
	}

	if err := app.container.Provide(repository.NewLessonRepository); err != nil {
		return err
	}

	return nil
}

func (app *Application) initServices() error {
	if err := app.container.Provide(func(repo repository.UserRepository, logger *zap.Logger) service.UserService {
		return service.NewUserService(repo, app.config.JWTSecret, app.config.JWTAccessExpiration, logger)
	}); err != nil {
		return err
	}

	if err := app.container.Provide(func(repo repository.LessonRepository, logger *zap.Logger) service.LessonService {
		return service.NewLessonService(repo, logger)
	}); err != nil {
		return err
	}

	return nil
}

func (app *Application) initControllers() error {
	if err := app.container.Invoke(func(
		e *echo.Echo,
		userService service.UserService,
		lessonService service.LessonService,
		logger *zap.Logger,
	) {
		// Swagger routes
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.GET("/swagger/doc.json", func(c echo.Context) error {
			return c.File("./docs/swagger.json")
		})
		e.GET("/swagger/doc.yaml", func(c echo.Context) error {
			return c.File("./docs/swagger.yaml")
		})

		// Пользовательские маршруты, доступные без JWT аутентификации
		userHandler := controller.NewUserController(userService, logger)
		e.POST("/api/v1/user/login", userHandler.Login)
		e.POST("/api/v1/user/register", userHandler.Register)

		// Маршруты уроков, доступные без JWT аутентификации
		lessonHandler := controller.NewLessonController(lessonService, logger)
		e.GET("/api/v1/lessons/types", lessonHandler.GetLessonTypes)
		e.GET("/api/v1/lessons/type/:type", lessonHandler.GetLessonsByType)
		e.GET("/api/v1/lessons/:id", lessonHandler.GetLesson)
		e.GET("/api/v1/lessons", lessonHandler.GetAllLessons)

		// Создаем группу для защищенных маршрутов
		api := e.Group("/api/v1")
		api.Use(middleware.JWTMiddleware(app.config.JWTSecret, app.config.JWTAccessExpiration, logger))

		// Защищенные маршруты пользователя
		api.GET("/user", userHandler.Profile)
	}); err != nil {
		return err
	}

	return nil
}

func (app *Application) Start() error {
	// Загружаем уроки из JSON файла при старте приложения
	return app.container.Invoke(func(e *echo.Echo, lessonService service.LessonService, logger *zap.Logger) error {
		// Проверяем существование файла с уроками
		lessonsFilePath := "data/lessons.json"
		if _, err := os.Stat(lessonsFilePath); err == nil {
			// Файл существует, загружаем уроки
			logger.Info("Начинаем загрузку уроков из файла", zap.String("path", lessonsFilePath))
			if err := lessonService.LoadLessonsFromFile(lessonsFilePath); err != nil {
				logger.Fatal("Ошибка загрузки уроков из файла", zap.Error(err))
			}
		} else {
			logger.Fatal("Файл с уроками не найден", zap.String("path", lessonsFilePath))
		}

		return e.Start(":" + app.config.Port)
	})
}

// GetContainer возвращает DI контейнер для тестирования
func (app *Application) GetContainer() *dig.Container {
	return app.container
}
