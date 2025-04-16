package config

import (
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
)

type AppConfig struct {
	Env     object.Env    `mapstructure:"env"`
	Logging LoggingConfig `mapstructure:"logging"`
}
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type GrpcServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type (
	DatabaseConfig struct {
		DataSource DataSourceConfig `mapstructure:"datasource"`
		Gorm       GormConfig       `mapstructure:"gorm"`
	}
	DataSourceConfig struct {
		Dsn string `mapstructure:"dsn"`
	}

	GormConfig struct {
		Logging LoggingConfig `mapstructure:"logging"`
	}
)

type (
	ConductorConfig struct {
		Client ConductorClientConfig `mapstructure:"client"`
	}
	ConductorClientConfig struct {
		Url         string `mapstructure:"root_uri"`
		ThreadCount int    `mapstructure:"thread_count"`
	}
)

type (
	TelemetryConfig struct {
		Collector CollectorConfig `mapstructure:"collector"`
	}
	CollectorConfig struct {
		Url string `mapstructure:"url"`
	}
)

type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	IsDebug bool   `mapstructure:"-"`
}

type KeycloakConfig struct {
	Url                string `mapstructure:"url"`
	Realm              string `mapstructure:"realm"`
	ClientId           string `mapstructure:"client_id"`
	ClientSecret       string `mapstructure:"client_secret"`
	RealmAdminUsername string `mapstructure:"realm_admin_username"`
	RealmAdminPassword string `mapstructure:"realm_admin_password"`
}
