package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	DatabaseDSN         string `mapstructure:"DATABASE_DSN"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	JWTAccessExpiration int    `mapstructure:"JWT_ACCESS_EXPIRATION"` // в часах
	Port                string `mapstructure:"PORT"`
}

func Load() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	envFile := ".env." + env

	viper.SetConfigName(envFile)
	viper.SetConfigType("env")

	viper.AddConfigPath("/var")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")

	// Устанавливаем значения по умолчанию
	viper.SetDefault("JWT_ACCESS_EXPIRATION", 24) // 24 часа по умолчанию

	viper.AutomaticEnv()

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Warning: failed to read config file %s: %v\n", envFile, err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Warning: failed to Unmarshal config file %s: %v\n", envFile, err)
	}

	return &cfg, nil
}
