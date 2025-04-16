//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/conductor"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/mysql"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/protostub"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/protosvc"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/web"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support"
	pkgconfig "github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
)

func InitializeServer(env object.Env, metadata *pkgconfig.Metadata) *infrastructure.Infra {
	wire.Build(
		config.NewConfig,
		AppConfigProviderSet,
		ServerConfigProviderSet,
		DatabaseConfigProviderSet,
		ConductorConfigProviderSet,
		TelemetryConfigProviderSet,
		GrpcConfigProviderSet,

		logging.NewLogHelper,

		database.NewDataSource,
		mysql.NewAuthorRepository,
		mysql.NewBookRepository,
		mysql.NewPublisherRepository,

		handler.NewAuthorHandler,
		handler.NewBookHandler,
		handler.NewCategoryHandler,
		handler.NewPublisherHandler,

		service.NewBookService,

		protostub.NewStockGrpcClient,
		protosvc.NewAuthorGrpcServer,
		protosvc.NewBookGrpcServer,
		protosvc.NewCategoryGrpcServer,
		support.NewGrpcServer,

		saga.NewWorkerClient,
		conductor.NewBookWorker,
		support.NewConductorWorker,

		web.NewAuthorController,
		web.NewBookController,
		web.NewCategoryController,
		web.NewPublisherController,
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
var ConductorConfigProviderSet = wire.NewSet(config.ProvideConductorConfig)
var TelemetryConfigProviderSet = wire.NewSet(config.ProvideTelemetryConfig)
var GrpcConfigProviderSet = wire.NewSet(config.ProvideGrpcConfig)
