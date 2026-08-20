package di

import (
	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/logging"
	"github.com/codejsha/shared-library-go/pkg/rest/client"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/workflow"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/pgsql"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/restcontroller"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support"
)

var Module = fx.Module("inventory",
	fx.Provide(
		config.NewConfig,
		config.ProvideAppConfig,
		config.ProvideServerConfig,
		config.ProvideDatabaseConfig,
		config.ProvideTelemetryConfig,
		config.ProvideTemporalConfig,
		config.ProvideKafkaConfig,

		pkgconfig.NewCloudConfigHelper,
		client.NewRestyClient,
		logging.NewLogHelper,

		restcontroller.NewStockController,
		restcontroller.NewWarehouseController,
		restcontroller.NewTransferController,
		restcontroller.NewAuditController,
		restcontroller.NewClosingController,
		restcontroller.NewBalanceController,

		service.NewStockService,

		database.NewVaultAwareDataSource,
		pgsql.NewStockRepository,
		pgsql.NewWarehouseRepository,
		pgsql.NewStockHistoryRepository,
		pgsql.NewStockReservationRepository,
		pgsql.NewStockTransferRepository,
		pgsql.NewStockAuditRepository,
		pgsql.NewMonthlyClosingRepository,
		pgsql.NewStockBalanceRepository,

		workflow.NewStockActivities,
		support.NewTemporalWorker,

		support.NewGinServer,
		support.NewTelemetryManager,
		support.NewKafkaPublisher,
		infrastructure.NewInfra,
	),
	fx.Invoke(support.RegisterEventPublisher),
)

func NewApp(preConfig *pkgconfig.PreConfig, metadata *pkgconfig.Metadata) *fx.App {
	return fx.New(
		fx.Supply(preConfig, metadata),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
}
