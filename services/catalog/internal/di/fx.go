package di

import (
	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/logging"
	"github.com/codejsha/shared-library-go/pkg/rest/client"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/cache"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/opensearch"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/pgsql"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support"
)

var Module = fx.Module("catalog",
	fx.Provide(
		config.NewConfig,
		config.ProvideAppConfig,
		config.ProvideServerConfig,
		config.ProvideDatabaseConfig,
		config.ProvideTelemetryConfig,
		config.ProvideOpensearchConfig,
		config.ProvideKafkaConfig,
		config.ProvideCacheConfig,

		pkgconfig.NewCloudConfigHelper,
		client.NewRestyClient,
		logging.NewLogHelper,

		restcontroller.NewWorkController,
		restcontroller.NewEditionController,
		restcontroller.NewAuthorController,
		restcontroller.NewPublisherController,
		restcontroller.NewSubjectController,

		service.NewCatalogService,

		database.NewVaultAwareDataSource,
		pgsql.NewWorkRepository,
		pgsql.NewEditionRepository,
		pgsql.NewAuthorRepository,
		pgsql.NewPublisherRepository,
		pgsql.NewSubjectRepository,

		opensearch.NewWorkSearchRepository,

		support.NewOpensearchClient,
		support.NewWorkListCache,

		support.NewGinServer,
		support.NewTelemetryManager,
		support.NewKafkaPublisher,
		infrastructure.NewInfra,
	),
	fx.Decorate(cache.NewCachingWorkRepo),
	fx.Invoke(support.RegisterEventPublisher),
	fx.Invoke(support.RegisterCacheInvalidator),
)

func NewApp(preConfig *pkgconfig.PreConfig, metadata *pkgconfig.Metadata) *fx.App {
	return fx.New(
		fx.Supply(preConfig, metadata),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
}
