package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	OtpTimeout                int // seconds
	ContextTimeout            int
	Environment               string // develop, staging, production
	LogLevel                  string // DEBUG, INFO ...
	HTTPPort                  string
	PostgresHost              string
	PostgresPort              string
	PostgresDatabase          string
	PostgresUser              string
	PostgresPassword          string
	PostgresConnectionTimeOut int // seconds
	PostgresConnectionTry     int
	BaseUrl                   string
	SignInKey                 string
	AuthConfigPath            string
	CSVFilePath               string
	RedisHost                 string
	RedisPort                 string
	AccessTokenTimout         int // MINUTES
	MaxImageSize              int // Mb
	RabbitMqHost              string
	RabbitMqPort              string
	RabbitMqUser              string
	RabbitMqPassword          string
	ListenQueue               string
	WriteQueue                string
}

// Load loads environment vars and inflates Config
func Load() Config {
	dotenvFilePath := cast.ToString(getOrReturnDefault("DOT_ENV_PATH", "/app/.env"))
	err := godotenv.Load(dotenvFilePath)

	if err != nil {
		fmt.Println(".env file not found. Default configuration is being used.")
	}
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "DEBUG"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", "8000"))
	c.BaseUrl = cast.ToString(getOrReturnDefault("BASE_URL", "http://localhost:8000/v1/"))

	// Postgres
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "transcoder"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "transcoder"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "secret"))
	c.PostgresConnectionTimeOut = cast.ToInt(getOrReturnDefault("POSTGRES_CONNECTION_TIMEOUT", 5))
	c.PostgresConnectionTry = cast.ToInt(getOrReturnDefault("POSTGRES_CONNECTION_TRY", 10))

	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", ""))
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/auth.csv"))

	// in mermory storage
	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))
	c.OtpTimeout = cast.ToInt(getOrReturnDefault("OTP_TIMEOUT", 300))
	c.ContextTimeout = cast.ToInt(getOrReturnDefault("CONTEXT_TIMOUT", 7))
	c.AccessTokenTimout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 300))

	// RabbitMQ
	c.RabbitMqHost = cast.ToString(getOrReturnDefault("RABBITMQ_HOST", "localhost"))
	c.RabbitMqPort = cast.ToString(getOrReturnDefault("RABBITMQ_PORT", "5672"))
	c.RabbitMqUser = cast.ToString(getOrReturnDefault("RABBITMQ_USER", "username"))
	c.RabbitMqPassword = cast.ToString(getOrReturnDefault("RABBITMQ_PASSWORD", "secret"))
	c.ListenQueue = cast.ToString(getOrReturnDefault("LISTEN_QUEUE", "pipeline_status"))
	c.WriteQueue = cast.ToString(getOrReturnDefault("WRITE_QUEUE", "pipelines"))

	// Media
	c.MaxImageSize = cast.ToInt(getOrReturnDefault("MAX_IMAGE_SIZE", 5))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
