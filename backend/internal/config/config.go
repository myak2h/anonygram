package config

import (
	"os"
	"strconv"

	"anonygram/internal/utils"
)

type Config struct {
	UploadPath       string
	MaxUploadSize    int64
	AllowedOrigins   []string
	Port             string
	ClientBufferSize int
	HubBufferSize    int
}

func Load() *Config {
	return &Config{
		UploadPath:       getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize:    getEnvInt64("MAX_UPLOAD_SIZE", 10<<20), // default 10 MB
		AllowedOrigins:   getEnvSlice("ALLOWED_ORIGINS", "*"),
		Port:             getEnv("PORT", "8080"),
		ClientBufferSize: getEnvInt("CLIENT_BUFFER_SIZE", 256),
		HubBufferSize:    getEnvInt("HUB_BUFFER_SIZE", 16),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if v := os.Getenv(key); v != "" {
		if int64Value, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int64Value
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvSlice(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	return utils.SplitAndTrim(value, ",")
}
