//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	pkgconfig "github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/mysql"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/protosvc"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/web"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support"
)

func InitializeServer(env object.Env, metadata *pkgconfig.Metadata) *infrastructure.Infra {
	wire.Build(
		config.NewConfig,
		AppConfigProviderSet,
		ServerConfigProviderSet,
		DatabaseConfigProviderSet,
		TelemetryConfigProviderSet,
		GrpcConfigProviderSet,

		logging.NewLogHelper,

		database.NewDataSource,
		mysql.NewStockRepository,
		mysql.NewWarehouseRepository,

		handler.NewStockHandler,
		handler.NewWarehouseHandler,

		service.NewStockService,

		protosvc.NewStockGrpcServer,
		support.NewGrpcServer,

		web.NewStockController,
		web.NewWarehouseController,
		support.NewGinServer,

		support.NewTelemetryManager,

		infrastructure.NewInfra,
	)
	return &infrastructure.Infra{}
}

func InitializeMigrationManager(env object.Env) *database.MigrationManager {
	wire.Build(
		config.NewConfig,
		DatabaseConfigProviderSet,
		database.NewMigrationManager,
	)
	return &database.MigrationManager{}
}

var AppConfigProviderSet = wire.NewSet(config.ProvideAppConfig)
var ServerConfigProviderSet = wire.NewSet(config.ProvideServerConfig)
var DatabaseConfigProviderSet = wire.NewSet(config.ProvideDatabaseConfig)
var TelemetryConfigProviderSet = wire.NewSet(config.ProvideTelemetryConfig)
var GrpcConfigProviderSet = wire.NewSet(config.ProvideGrpcConfig)
