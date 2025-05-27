package main

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	_ "github.com/Ravr-Site/Ravr-Backend/docs" // Импорт для Swagger документации
	"github.com/Ravr-Site/Ravr-Backend/internal"
	"log"
)

// @title RAVR Backend API
// @version 1.0
// @description API для проекта RAVR
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.ravr-site.io/support
// @contact.email support@ravr-site.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := internal.NewApplication(cfg)
	if err := app.Init(); err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}
}
