package internal

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/controller"
	"github.com/Ravr-Site/Ravr-Backend/internal/database"
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
	"path/filepath"
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

		// Запускаем миграции с помощью golang-migrate
		if err := app.runMigrations(db, logger); err != nil {
			logger.Error("Failed to run migrations", zap.Error(err))
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

// runMigrations выполняет миграции базы данных с помощью golang-migrate
func (app *Application) runMigrations(db *gorm.DB, logger *zap.Logger) error {
	logger.Info("Running database migrations")

	// Определяем путь к директории с миграциями
	// Сначала пытаемся найти относительно текущего рабочего каталога
	migrationsPath := "migrations"

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Если директория не найдена, используем абсолютный путь от корня проекта
		workDir, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get working directory", zap.Error(err))
			return err
		}

		// Пробуем найти директорию миграций относительно корня проекта
		migrationsPath = filepath.Join(workDir, "migrations")
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			logger.Error("Migrations directory not found", zap.String("path", migrationsPath))
			return err
		}
	}

	logger.Info("Using migrations path", zap.String("path", migrationsPath))

	// Проверяем переменную окружения для принудительного сброса миграций
	if os.Getenv("RESET_MIGRATIONS") == "true" {
		logger.Warn("Reset migrations flag is set, reverting all migrations")
		if err := database.ResetMigrations(db, migrationsPath); err != nil {
			logger.Error("Failed to reset migrations", zap.Error(err))
			return err
		}
	}

	// Запускаем миграции
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return err
	}

	logger.Info("Database migrations completed successfully")
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
