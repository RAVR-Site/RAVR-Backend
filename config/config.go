package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseDSN string
	JWTSecret   string
}

func Load() (*Config, error) {
	dsn := os.Getenv("DATABASE_DSN")
	jwt := os.Getenv("JWT_SECRET")
	return &Config{DatabaseDSN: dsn, JWTSecret: jwt}, nil
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
