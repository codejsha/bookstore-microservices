package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/codejsha/shared-library-go/pkg/config"
)

const (
	vaultSecretsFile  = "/vault/secrets/db.properties"
	configFileEnv     = "APP_CONFIG_FILE"
	defaultConfigFile = "/app/config/application.yaml"
)

type Config struct {
	App       *config.AppConfig       `mapstructure:"app"`
	Server    *config.ServerConfig    `mapstructure:"server"`
	Database  *config.DatabaseConfig  `mapstructure:"database"`
	Telemetry *config.TelemetryConfig `mapstructure:"telemetry"`
	Kafka     *config.KafkaConfig     `mapstructure:"kafka"`
	Grpc      *GrpcConfig             `mapstructure:"grpc"`
}

type GrpcConfig struct {
	OrderServer       *config.GrpcServerConfig `mapstructure:"order_server"`
	PaymentServer     *config.GrpcServerConfig `mapstructure:"payment_server"`
	UserServer        *config.GrpcServerConfig `mapstructure:"user_server"`
	FulfillmentServer *config.GrpcServerConfig `mapstructure:"fulfillment_server"`
}

func NewConfig(
	preConfig *config.PreConfig,
) *Config {
	path := configFilePath()
	cfg, err := readConfigFile(path)
	if err == nil {
		config.ApplyVaultDBCredentials(cfg.Database, vaultSecretsFile)
		return cfg
	}
	if errors.Is(err, fs.ErrNotExist) {
		logrus.Infof("config file %s not found, reading local config files", path)
	} else {
		logrus.Errorf("failed to read config file %s: %v", path, err)
	}

	cfg, err = readConfig(preConfig.Profile)
	if err == nil {
		config.ApplyVaultDBCredentials(cfg.Database, vaultSecretsFile)
		return cfg
	}
	logrus.Errorf("failed to read config from local config files: %v", err)

	panic("failed to load config")
}

func configFilePath() string {
	if path := os.Getenv(configFileEnv); path != "" {
		return path
	}
	return defaultConfigFile
}

func readConfigFile(
	path string,
) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	vp := viper.New()
	vp.SetConfigType("yaml")
	if err := vp.ReadConfig(bytes.NewReader(content)); err != nil {
		return nil, fmt.Errorf("fatal error reading config file %s: %w", path, err)
	}
	for _, key := range vp.AllKeys() {
		vp.Set(key, expandEnvValues(vp.Get(key)))
	}

	return decodeConfig(vp)
}

func expandEnvValues(value any) any {
	switch v := value.(type) {
	case string:
		return os.ExpandEnv(v)
	case []any:
		expanded := make([]any, len(v))
		for i, item := range v {
			expanded[i] = expandEnvValues(item)
		}
		return expanded
	case []string:
		expanded := make([]string, len(v))
		for i, item := range v {
			expanded[i] = os.ExpandEnv(item)
		}
		return expanded
	case map[string]any:
		expanded := make(map[string]any, len(v))
		for k, item := range v {
			expanded[k] = expandEnvValues(item)
		}
		return expanded
	default:
		return v
	}
}

func readConfig(
	profile config.Profile,
) (*Config, error) {
	vp := viper.New()
	vp.SetConfigName(fmt.Sprintf("config-%s", profile))
	vp.SetConfigType("json")
	vp.AddConfigPath("configs")

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %s\n", err)
	}

	return decodeConfig(vp)
}

func decodeConfig(
	vp *viper.Viper,
) (*Config, error) {
	cfg := &Config{}
	if err := vp.Unmarshal(cfg, viper.DecodeHook(config.RejectScalarSliceHook())); err != nil {
		return nil, fmt.Errorf("fatal error unmarshaling config: %w", err)
	}
	if cfg.App == nil {
		return nil, errors.New("config is missing the app section")
	}
	cfg.App.Logging.UpdateFlags()
	return cfg, nil
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
func ProvideTelemetryConfig(cfg *Config) *config.TelemetryConfig {
	return cfg.Telemetry
}
func ProvideKafkaConfig(cfg *Config) *config.KafkaConfig {
	return cfg.Kafka
}
func ProvideGrpcConfig(cfg *Config) *GrpcConfig {
	return cfg.Grpc
}
