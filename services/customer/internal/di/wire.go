//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	pkgconfig "github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/conductor"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/mysql"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/web"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
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
		mysql.NewPointRepository,
		mysql.NewWishlistRepository,

		mapper.NewOrderMapper,
		mapper.NewPaymentMapper,
		mapper.NewUserMapper,

		handler.NewPointHandler,
		handler.NewWishlistHandler,

		service.NewCustomerService,
		service.NewOrderService,
		service.NewPaymentService,
		service.NewPointService,
		service.NewWishlistService,

		protostub.NewOrderGrpcClient,
		protostub.NewPaymentGrpcClient,
		protostub.NewUserGrpcClient,

		saga.NewWorkerClient,
		conductor.NewCustomerWorker,
		support.NewConductorWorker,

		web.NewCustomerController,
		web.NewOrderController,
		web.NewPaymentController,
		web.NewPointController,
		web.NewWishlistController,
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
