package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	DatabaseDSN string `mapstructure:"DATABASE_DSN"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Port        string `mapstructure:"PORT"`
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

	viper.AutomaticEnv()

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Warning: failed to read config file %s: %v\n", envFile, err)
	}

	dbHost := viper.GetString("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dsnFormat := viper.GetString("DATABASE_DSN")
	formattedDSN := fmt.Sprintf(dsnFormat, dbHost)
	viper.Set("DATABASE_DSN", formattedDSN)

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Warning: failed to Unmarshal config file %s: %v\n", envFile, err)
	}

	return &cfg, nil
}
