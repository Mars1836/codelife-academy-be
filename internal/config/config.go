package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address         string
	DatabaseURL     string
	RedisAddress    string
	RedisPassword   string
	RedisDB         int
	CacheTTL        time.Duration
	ShutdownTimeout time.Duration
	MaxBodyBytes    int64
	AuthTokenSecret string
	AuthTokenTTL    time.Duration
	AuthOTPTTL      time.Duration
	SMTPHost        string
	SMTPPort        string
	SMTPUsername    string
	SMTPPassword    string
	SMTPFrom        string
}

func Load() Config {
	return Config{
		Address:         getenv("HTTP_ADDRESS", ":8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		RedisAddress:    os.Getenv("REDIS_ADDRESS"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:         getint("REDIS_DB", 0),
		CacheTTL:        time.Duration(getint("CACHE_TTL_SECONDS", 300)) * time.Second,
		ShutdownTimeout: time.Duration(getint("SHUTDOWN_TIMEOUT_SECONDS", 10)) * time.Second,
		MaxBodyBytes:    int64(getint("MAX_BODY_BYTES", 1<<20)),
		AuthTokenSecret: os.Getenv("AUTH_TOKEN_SECRET"),
		AuthTokenTTL:    time.Duration(getint("AUTH_TOKEN_TTL_SECONDS", 86400)) * time.Second,
		AuthOTPTTL:      time.Duration(getint("AUTH_OTP_TTL_SECONDS", 600)) * time.Second,
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        getenv("SMTP_PORT", "587"),
		SMTPUsername:    os.Getenv("SMTP_USERNAME"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        os.Getenv("SMTP_FROM"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getint(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
