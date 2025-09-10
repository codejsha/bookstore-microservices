package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/codejsha/shared-library-go/pkg/config"
)

const vaultSecretsFile = "/vault/secrets/db.properties"

type Config struct {
	App        *config.AppConfig       `mapstructure:"app"`
	Server     *config.ServerConfig    `mapstructure:"server"`
	Database   *config.DatabaseConfig  `mapstructure:"database"`
	Cache      *config.CacheConfig     `mapstructure:"cache"`
	Telemetry  *config.TelemetryConfig `mapstructure:"telemetry"`
	Kafka      *config.KafkaConfig     `mapstructure:"kafka"`
	Opensearch *OpensearchConfig       `mapstructure:"opensearch"`
}

type OpensearchConfig struct {
	Scheme   string `mapstructure:"scheme"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Index    string `mapstructure:"index"`
	Insecure bool   `mapstructure:"insecure"`
}

func NewConfig(
	preConfig *config.PreConfig,
	cloudConfigHelper *config.CloudConfigHelper,
) *Config {
	cfg, err := fetchConfig(*preConfig, *cloudConfigHelper)
	if err == nil {
		config.ApplyVaultDBCredentials(cfg.Database, vaultSecretsFile)
		return cfg
	}
	logrus.Errorf("failed to fetch config from config service: %v", err)

	cfg, err = readConfig(preConfig.Profile)
	if err == nil {
		config.ApplyVaultDBCredentials(cfg.Database, vaultSecretsFile)
		return cfg
	}
	logrus.Errorf("failed to read config from local config files: %v", err)

	panic("failed to load config")
}

func fetchConfig(
	preConfig config.PreConfig,
	cloudConfigHelper config.CloudConfigHelper,
) (*Config, error) {
	cfgResp, err := cloudConfigHelper.FetchConfig(
		preConfig.ConfigServerAddr,
		preConfig.ServiceName,
		string(preConfig.Profile),
		preConfig.Label,
	)
	if err != nil {
		return nil, err
	}

	source := config.NormalizePropertySources(cfgResp.PropertySources)
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal flattened config: %w", err)
	}

	vp := viper.New()
	vp.SetConfigType("json")
	if err := vp.ReadConfig(bytes.NewBufferString(os.ExpandEnv(string(sourceBytes)))); err != nil {
		return nil, fmt.Errorf("fatal error reading config data: %w", err)
	}

	cfg := &Config{}
	if err := vp.Unmarshal(cfg, viper.DecodeHook(config.RejectScalarSliceHook())); err != nil {
		return nil, fmt.Errorf("fatal error unmarshaling config data: %w", err)
	}
	cfg.App.Logging.UpdateFlags()
	return cfg, nil
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

	cfg := &Config{}
	if err := vp.Unmarshal(cfg, viper.DecodeHook(config.RejectScalarSliceHook())); err != nil {
		return nil, fmt.Errorf("fatal error unmarshaling config file: %s\n", err)
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
func ProvideCacheConfig(cfg *Config) *config.CacheConfig {
	return cfg.Cache
}
func ProvideTelemetryConfig(cfg *Config) *config.TelemetryConfig {
	return cfg.Telemetry
}
func ProvideKafkaConfig(cfg *Config) *config.KafkaConfig {
	return cfg.Kafka
}
func ProvideOpensearchConfig(cfg *Config) *OpensearchConfig {
	return cfg.Opensearch
}
