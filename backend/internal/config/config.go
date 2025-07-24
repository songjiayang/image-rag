package config

import (
	"log"

	"image-rag-backend/internal/utils"

	"github.com/joho/godotenv"
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
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file, using system environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     utils.GetEnv("MYSQL_HOST", "localhost"),
			Port:     utils.GetEnv("MYSQL_PORT", "3306"),
			User:     utils.GetEnv("MYSQL_USER", "root"),
			Password: utils.GetEnv("MYSQL_PASSWORD", "password"),
			Database: utils.GetEnv("MYSQL_NAME", "image_rag"),
		},
		Server: ServerConfig{
			Port: utils.GetEnv("SERVER_PORT", "8080"),
		},
		Upload: UploadConfig{
			Path:      utils.GetEnv("UPLOAD_PATH", "./uploads"),
			MaxSizeMB: int64(utils.GetEnvInt("MAX_UPLOAD_SIZE_MB", 10)),
			AllowedExt: map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".webp": true,
			},
		},
		Doubao: DoubaoConfig{
			APIKey: utils.GetEnv("DOUBAO_API_KEY", ""),
			Model:  utils.GetEnv("DOUBAO_MODEL", "doubao-embedding-vision-250615"),
			URL:    utils.GetEnv("DOUBAO_API_URL", "https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal"),
		},
		Milvus: MilvusConfig{
			Host:     utils.GetEnv("MILVUS_HOST", "localhost"),
			Port:     utils.GetEnv("MILVUS_PORT", "19530"),
			Database: utils.GetEnv("MILVUS_DATABASE", "image_rag"),
		},
	}
}
