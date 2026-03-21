package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	AppPort string
	AppEnv  string
}

type JWTConfig struct {
	Secret             string
	ExpireMinutes      int
	RefreshExpireHours int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file, reading from environment")
	}

	jwtExpireMinutes, _ := strconv.Atoi(getEnv("JWT_EXPIRE_MINUTES", "15"))
	jwtRefreshExpireHours, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_HOURS", "24"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "crm_admin"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			AppPort: getEnv("APP_PORT", "8080"),
			AppEnv:  getEnv("APP_ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			ExpireMinutes:      jwtExpireMinutes,
			RefreshExpireHours: jwtRefreshExpireHours,
		},
	}

}

func getEnv(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (s *ServerConfig) IsDevelopment() bool {
	return s.AppEnv == "development"
}
