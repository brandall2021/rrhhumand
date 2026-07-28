package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string
	AppDebug bool
	JWT     JWTConfig
	MinIO   MinIOConfig
	MaxUploadSizeMB int
}

type JWTConfig struct {
	Secret          string
	Expiration      time.Duration
	RefreshExpiration time.Duration
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

func Load() *Config {
	godotenv.Load(".env")

	jwtExp, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "15m"))
	jwtRefreshExp, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRATION", "7d"))
	maxUpload, _ := strconv.Atoi(getEnv("MAX_UPLOAD_SIZE_MB", "25"))
	if maxUpload <= 0 {
		maxUpload = 25
	}

	return &Config{
		AppName:  getEnv("APP_NAME", "rrhhumand"),
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("APP_PORT", "8080"),
		AppDebug: getEnv("APP_DEBUG", "true") == "true",
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "change-me-in-production"),
			Expiration:        jwtExp,
			RefreshExpiration: jwtRefreshExp,
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "rrhhumand"),
			Region:    getEnv("MINIO_REGION", "us-east-1"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		MaxUploadSizeMB: maxUpload,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
