package config

import "time"

type Config struct {
	DatabaseURL             string
	SecretKey               string
	Port                    string
	AccessTokenExpiresTime  time.Duration
	RefreshTokenExpiresTime time.Duration
}

var AppConfig = Config{
	DatabaseURL:             "postgres://postgres:Zbekxzz3$$@localhost:5432/auth?sslmode=disable",
	SecretKey:               "super_secret_key",
	Port:                    "3000",
	AccessTokenExpiresTime:  15 * time.Minute,
	RefreshTokenExpiresTime: 24 * time.Hour,
}
