package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL             string
	SecretKey               string
	Port                    string
	AccessTokenExpiresTime  time.Duration
	RefreshTokenExpiresTime time.Duration
}

var AppConfig = Config{
	DatabaseURL:             os.Getenv("DATABASE_URL"),
	SecretKey:               "explore",
	Port:                    "3000",
	AccessTokenExpiresTime:  15 * time.Minute,
	RefreshTokenExpiresTime: 24 * time.Hour,
}
