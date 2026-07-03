package config

import "os"

type Config struct {
	DatabaseURL string
	RedisHost   string
	RedisPort   string
	APIKey      string
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func LoadConfig() *Config {

	if getEnv("DATABASE_URL") == "" {
		panic("DATABASE_URL is required")
	}

	if getEnv("REDIS_HOST") == "" {
		panic("REDIS_HOST is required")
	}

	if getEnv("REDIS_PORT") == "" {
		panic("REDIS_PORT is required")
	}

	if getEnv("API_KEY") == "" {
		panic("API_KEY is required")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL"),
		RedisHost:   getEnv("REDIS_HOST"),
		RedisPort:   getEnv("REDIS_PORT"),
		APIKey:      getEnv("API_KEY"),
	}
}
