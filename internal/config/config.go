package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type Config struct {
	ServerPort string
	DBPath     string
	JWTSecret  string
	LogDir     string
}

func LoadConfig() (*Config, error) {
	// Ignore error as .env might not exist in production
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "myapp.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey" // Default for local dev
	}

	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}

	return &Config{
		ServerPort: port,
		DBPath:     dbPath,
		JWTSecret:  jwtSecret,
		LogDir:     logDir,
	}, nil
}

// Module provides the configuration to Fx
var Module = fx.Provide(LoadConfig)
