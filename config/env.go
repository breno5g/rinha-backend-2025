package config

import (
	"github.com/spf13/viper"
)

type conf struct {
	DefaultURL  string `mapstructure:"DEFAULT_URL"`
	FallbackURL string `mapstructure:"FALLBACK_URL"`
	MaxWorkers  int    `mapstructure:"MAX_WORKERS"`
	Port        int    `mapstructure:"PORT"`
	RedisURL    string `mapstructure:"REDIS_URL"`
}

func InitEnv() (*conf, error) {
	cfg := &conf{}

	viper.SetDefault("DEFAULT_URL", "http://localhost:8001/payments")
	viper.SetDefault("FALLBACK_URL", "http://localhost:8002/payments")
	viper.SetDefault("MAX_WORKERS", 5)
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("REDIS_URL", "localhost:6379")

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
