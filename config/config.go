package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   ServerConfig
		Postgres PostgresConfig
		Redis    RedisConfig
		JWT      JWTConfig
		Email    EmailConfig
		Logger   LoggerConfig
	}

	JWTConfig struct {
		SecretKey string
	}

	ServerConfig struct {
		Port string
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	RedisConfig struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	LoggerConfig struct {
		Level string
	}

	EmailConfig struct {
		SmtpHost string
		SmtpPort int
		SmtpUser string
		SmtpPass string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	// Server
	c.Server.Port = os.Getenv("SERVER_PORT")

	// Postgres
	c.Postgres.Host = os.Getenv("DB_HOST")
	c.Postgres.Port = os.Getenv("DB_PORT")
	c.Postgres.User = os.Getenv("DB_USER")
	c.Postgres.Password = os.Getenv("DB_PASSWORD")
	c.Postgres.DBName = os.Getenv("DB_NAME")

	// Redis
	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")
	c.Redis.Username = os.Getenv("REDIS_USERNAME")
	c.Redis.Password = os.Getenv("REDIS_PASSWORD")

	// Logger
	c.Logger.Level = os.Getenv("LOG_LEVEL")

	// Email
	c.Email.SmtpHost = os.Getenv("SMTP_HOST")
	c.Email.SmtpPort = smtpPort
	c.Email.SmtpUser = os.Getenv("SMTP_USER")
	c.Email.SmtpPass = os.Getenv("SMTP_PASS")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}
