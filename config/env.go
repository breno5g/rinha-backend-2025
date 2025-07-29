package config

import (
	"github.com/spf13/viper"
)

type conf struct {
	DefaultURL  string `mapstructure:"DEFAULT_URL"`
	FallbackURL string `mapstructure:"FALLBACK_URL"`
	MaxWorkers  int    `mapstructure:"MAX_WORKERS"`
}

func InitEnv(path string) (*conf, error) {
	cfg := &conf{
		DefaultURL:  "http://localhost:8001/payments",
		FallbackURL: "http://localhost:8002/payments",
		MaxWorkers:  5,
	}

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
