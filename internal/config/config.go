package config

import (
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Metrics   MetricsConfig   `mapstructure:"metrics" validate:"required"`
		Analyzers AnalyzersConfig `mapstructure:"analyzers" validate:"required"`
	}

	MetricsConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port" validate:"required"`
	}

	AnalyzersConfig struct {
		KRR      KrrAnalyzerConfig `mapstructure:"krr" validate:"required"`
		VPC      VPCAnalyzerConfig `mapstructure:"vpc" validate:"required"`
		CloudID  string            `mapstructure:"cloud_id"`
		FolderID string            `mapstructure:"folder_id"`
	}

	KrrAnalyzerConfig struct {
		PrometheusURL        string `mapstructure:"prometheus_url" validate:"required"`
		PrometheusAuthHeader string `mapstructure:"prometheus_auth_header" validate:"required"`
		HistoryDuration      string `mapstructure:"history_duration" validate:"required"`
	}

	VPCAnalyzerConfig struct {
		YCToken string `mapstructure:"yc_token" validate:"required"`
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
