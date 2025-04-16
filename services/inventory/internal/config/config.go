package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
)

type Config struct {
	App       *config.AppConfig       `mapstructure:"app"`
	Server    *config.ServerConfig    `mapstructure:"server"`
	Database  *config.DatabaseConfig  `mapstructure:"database"`
	Conductor *config.ConductorConfig `mapstructure:"conductor"`
	Telemetry *config.TelemetryConfig `mapstructure:"telemetry"`
	Grpc      *GrpcConfig             `mapstructure:"grpc"`
}

type GrpcConfig struct {
	Server *config.ServerConfig `mapstructure:"server"`
}

func NewConfig(env object.Env) *Config {
	viper.SetConfigName(fmt.Sprintf("config-%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error config file: %s\n", err)
		return nil
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		_ = fmt.Errorf("fatal error unmarshaling config file: %s\n", err)
		return nil
	}
	updateConfigProperty(&cfg)
	return &cfg
}

func updateConfigProperty(cfg *Config) {
	switch cfg.App.Logging.Level {
	case "debug":
		cfg.App.Logging.IsDebug = true
	default:
		cfg.App.Logging.IsDebug = false
	}
}

func ProvideAppConfig(cfg *Config) *config.AppConfig {
	return cfg.App
}
func ProvideServerConfig(cfg *Config) *config.ServerConfig {
	return cfg.Server
}
func ProvideDatabaseConfig(cfg *Config) *config.DatabaseConfig {
	return cfg.Database
}
func ProvideConductorConfig(cfg *Config) *config.ConductorConfig {
	return cfg.Conductor
}
func ProvideTelemetryConfig(cfg *Config) *config.TelemetryConfig {
	return cfg.Telemetry
}
func ProvideGrpcConfig(cfg *Config) *GrpcConfig {
	return cfg.Grpc
}
