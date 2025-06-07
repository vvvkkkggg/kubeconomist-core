package config

import (
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		MetricsPort string `mapstructure:"metrics_port" validate:"required"`
	}
)

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, validatorPkg.New().Struct(&config)
}
