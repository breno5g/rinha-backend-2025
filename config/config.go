package config

import (
	"github.com/redis/go-redis/v9"
)

var (
	db     *redis.Client
	logger *Logger
	env    *conf
)

func Init() {
	var err error
	logger := NewLogger("Config")

	env, err = InitEnv(".")
	if err != nil {
		logger.Error("Error loading environment variables", err)
		panic(err)
	}

	db, err = NewRedisClient()

	if err != nil {
		logger.Error("Error connecting to Redis", err)
		panic(err)
	}
}

func GetDB() *redis.Client {
	return db
}

func GetLogger() *Logger {
	return logger
}

func GetEnv() *conf {
	return env
}
