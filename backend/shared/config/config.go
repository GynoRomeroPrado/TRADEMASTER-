package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment string
	Port        int
	Database    DatabaseConfig
	Redis       RedisConfig
	MongoDB     MongoDBConfig
	RabbitMQ    RabbitMQConfig
	JWT         JWTConfig
	AWS         AWSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	URL string
}

type MongoDBConfig struct {
	URI string
}

type RabbitMQConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
	Expiry string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3BucketPhotos  string
	S3BucketBlueprints string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("NODE_ENV", "development"),
		Port:        getEnvAsInt("PORT", 8080),
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "trademaster"),
			Password: getEnv("POSTGRES_PASSWORD", "changeme"),
			DBName:   getEnv("POSTGRES_DB", "trademaster"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		MongoDB: MongoDBConfig{
			URI: getEnv("MONGODB_URI", "mongodb://localhost:27017/trademaster"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		AWS: AWSConfig{
			Region:             getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:        getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:    getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3BucketPhotos:     getEnv("S3_BUCKET_PHOTOS", "trademaster-photos"),
			S3BucketBlueprints: getEnv("S3_BUCKET_BLUEPRINTS", "trademaster-blueprints"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
