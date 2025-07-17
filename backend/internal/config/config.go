package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Upload   UploadConfig
	Doubao   DoubaoConfig
	Milvus   MilvusConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type ServerConfig struct {
	Port string
}

type UploadConfig struct {
	Path       string
	MaxSizeMB  int64
	AllowedExt map[string]bool
}

type DoubaoConfig struct {
	APIKey string
	Model  string
	URL    string
}

type MilvusConfig struct {
	Host     string
	Port     string
	Database string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "image_rag"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Upload: UploadConfig{
			Path:      getEnv("UPLOAD_PATH", "./uploads"),
			MaxSizeMB: int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 10)),
			AllowedExt: map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".webp": true,
			},
		},
		Doubao: DoubaoConfig{
			APIKey: getEnv("DOUBAO_API_KEY", ""),
			Model:  getEnv("DOUBAO_MODEL", "doubao-embedding-vision-250615"),
			URL:    getEnv("DOUBAO_API_URL", "https://ark.cn-beijing.volces.com/api/v3/embeddings"),
		},
		Milvus: MilvusConfig{
			Host:     getEnv("MILVUS_HOST", "localhost"),
			Port:     getEnv("MILVUS_PORT", "19530"),
			Database: getEnv("MILVUS_DATABASE", "image_rag"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
