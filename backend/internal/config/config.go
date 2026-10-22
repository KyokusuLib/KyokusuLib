package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultJWTSecret   = "jwt_secret"
	defaultDatabaseURL = "postgres://user:pass@localhost:5432/mydb?sslmode=disable"

	minJWTSecretLength = 32
)

type Config struct {
	Address      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	GoogleClientID          string
	GoogleClientSecret      string
	GoogleRedirectURL       string
	GoogleSocialRedirectURL string

	DiscordClientID          string
	DiscordClientSecret      string
	DiscordRedirectURL       string
	DiscordSocialRedirectURL string

	JWTSecret string

	DatabaseURL string
	FrontendURL string
	BackendURL  string

	CookieSecure bool

	KyokusuEmailPass string
	KyokusuEmailName string

	ParserURL string

	RedisAddr     string
	RedisPassword string

	TgBotApiKey string

	TgBotToken        string
	TgChannelUsername string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		Address:      getEnv("SERVER_ADDRESS", "127.0.0.1:8080"),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),

		GoogleClientID:          getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:       getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		GoogleSocialRedirectURL: getEnv("SOCIAL_GOOGLE_REDIRECT_URL", "http://localhost:8080/api/socials/google/callback"),

		DiscordClientID:          getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret:      getEnv("DISCORD_CLIENT_SECRET", ""),
		DiscordRedirectURL:       getEnv("DISCORD_REDIRECT_URL", "http://localhost:8080/api/auth/discord/callback"),
		DiscordSocialRedirectURL: getEnv("SOCIAL_DISCORD_CLIENT_URL", "http://localhost:8080/api/socials/discord/callback"),

		DatabaseURL: getEnv("DATABASE_URL", defaultDatabaseURL),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		BackendURL:  getEnv("BACKEND_URL", "http://localhost:8080"),

		CookieSecure: getEnvBool("COOKIE_SECURE", false),

		JWTSecret: getEnv("JWT_SECRET", defaultJWTSecret),

		KyokusuEmailPass: getEnv("KYOKUSU_EMAIL_PASSWORD", ""),
		KyokusuEmailName: getEnv("KYOKUSU_EMAIL_NAME", ""),

		ParserURL: getEnv("PARSER_URL", "http://localhost:3005"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		TgBotApiKey: getEnv("TG_BOT_API_KEY", ""),

		TgBotToken:        getEnv("TELEGRAM_BOT_TOKEN", ""),
		TgChannelUsername: getEnv("TG_CHANNEL_USERNAME", "KyokusuLib"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var errs []error

	switch {
	case c.JWTSecret == "" || c.JWTSecret == defaultJWTSecret:
		errs = append(errs, errors.New("JWT_SECRET must be set to a non-default value"))
	case len(c.JWTSecret) < minJWTSecretLength:
		errs = append(errs, fmt.Errorf("JWT_SECRET must be at least %d characters long", minJWTSecretLength))
	}

	if c.DatabaseURL == "" || c.DatabaseURL == defaultDatabaseURL {
		errs = append(errs, errors.New("DATABASE_URL must be set"))
	}

	return errors.Join(errs...)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration for %s: %s. Using default %v", key, value, defaultVal)
		return defaultVal
	}
	return duration
}

func getEnvBool(key string, defaultVal bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		log.Printf("Invalid bool for %s: %q. Using default %v", key, value, defaultVal)
		return defaultVal
	}
}
