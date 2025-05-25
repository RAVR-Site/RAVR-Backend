package main

import (
	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal"
	"log"
)

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
