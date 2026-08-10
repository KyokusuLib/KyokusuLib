package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Address      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleSocialRedirectURL  string

	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURL  string
	DiscordSocialRedirectURL  string

	JWTSecret   string
	StateString string

	DatabaseURL string
	FrontendURL string
	BackendURL  string

	KyokusuEmailPass string
	KyokusuEmailName string

	ParserURL string

	RedisAddr     string
	RedisPassword string

	TgBotApiKey string

	TgBotToken        string
	TgChannelUsername string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Address:      getEnv("SERVER_ADDRESS", "127.0.0.1:8080"),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		GoogleSocialRedirectURL:  getEnv("SOCIAL_GOOGLE_REDIRECT_URL", "http://localhost:8080/api/socials/google/callback"),

		DiscordClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		DiscordRedirectURL:  getEnv("DISCORD_REDIRECT_URL", "http://localhost:8080/api/auth/discord/callback"),
		DiscordSocialRedirectURL:  getEnv("SOCIAL_DISCORD_CLIENT_URL", "http://localhost:8080/api/socials/discord/callback"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/mydb?sslmode=disable"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		BackendURL:  getEnv("BACKEND_URL", "http://localhost:8080"),

		JWTSecret:   getEnv("JWT_SECRET", "jwt_secret"),
		StateString: getEnv("STATE", "random-string-secret"),

		KyokusuEmailPass: getEnv("KYOKUSU_EMAIL_PASSWORD", ""),
		KyokusuEmailName: getEnv("KYOKUSU_EMAIL_NAME", ""),

		ParserURL: getEnv("PARSER_URL", "http://localhost:3005"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		TgBotApiKey: getEnv("TG_BOT_API_KEY", ""),

		TgBotToken:        getEnv("TELEGRAM_BOT_TOKEN", ""),
		TgChannelUsername: getEnv("TG_CHANNEL_USERNAME", "KyokusuLib"),
	}
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
