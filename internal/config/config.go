package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBType     string
	ServerPort string
	AppEnv     string

	// PostgreSQL Configuration
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

// Load parses configuration from .env, flags, and environment variables.
func Load() *Config {
	cfg := &Config{}

	// Load .env file for local development.
	// We ignore the error so it doesn't crash if the file doesn't exist in staging/prod.
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] No .env file found, reading from system environment/flags")
	} else {
		log.Println("[Config] Successfully loaded configurations from local .env file")
	}

	// Define command-line flags (which override .env values if passed explicitly)
	flag.StringVar(&cfg.DBType, "db-type", getEnv("DB_TYPE", "memory"), "Database type: 'memory' or 'postgres'")
	flag.StringVar(&cfg.ServerPort, "port", getEnv("SERVER_PORT", "8080"), "HTTP server listen port")

	// PostgreSQL specific flags
	flag.StringVar(&cfg.DBHost, "db-host", getEnv("DB_HOST", "localhost"), "PostgreSQL server host/IP")
	flag.StringVar(&cfg.DBPort, "db-port", getEnv("DB_PORT", "5432"), "PostgreSQL server port")
	flag.StringVar(&cfg.DBName, "db-name", getEnv("DB_NAME", "todo_db"), "PostgreSQL database name")
	flag.StringVar(&cfg.DBUser, "db-user", getEnv("DB_USER", "postgres"), "PostgreSQL username")
	flag.StringVar(&cfg.DBPassword, "db-password", getEnv("DB_PASSWORD", ""), "PostgreSQL password")
	flag.StringVar(&cfg.DBSSLMode, "db-sslmode", getEnv("DB_SSLMODE", "disable"), "PostgreSQL SSL mode")
	flag.StringVar(&cfg.AppEnv, "app-env", getEnv("APP_ENV", "development"), "Application environment; swagger disabled in production")

	flag.Parse()

	return cfg
}

// GetDatabaseURL constructs the PostgreSQL connection string from individual components
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func (c *Config) IsSwaggerEnabled() bool {
	return strings.ToLower(c.AppEnv) != "production"
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
