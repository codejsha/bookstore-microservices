package di

import (
	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/logging"
	"github.com/codejsha/shared-library-go/pkg/rest/client"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/cache"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/pgsql"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/protosvc"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/restcontroller"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/support"
)

var Module = fx.Module("identity",
	fx.Provide(
		config.NewConfig,
		config.ProvideAppConfig,
		config.ProvideServerConfig,
		config.ProvideDatabaseConfig,
		config.ProvideTelemetryConfig,
		config.ProvideGrpcConfig,
		config.ProvideCacheConfig,
		config.ProvideKafkaConfig,

		pkgconfig.NewCloudConfigHelper,
		client.NewRestyClient,
		logging.NewLogHelper,

		restcontroller.NewUserController,
		restcontroller.NewAuthzController,
		restcontroller.NewRiskController,

		service.NewIdentityService,

		database.NewVaultAwareDataSource,
		pgsql.NewUserRepository,

		protosvc.NewUserGrpcServer,

		keycloak.NewAdminTokenHelper,
		keycloak.NewUsersClient,
		keycloak.NewIntrospector,

		support.NewCacheClient,
		support.NewAdminTokenCache,
		support.NewUserCacheStore,
		support.NewRevocationCache,
		support.NewRiskCache,
		support.ProvideRiskStore,
		support.ProvideRiskChecker,
		support.ProvideSessionRevoker,
		support.ProvideRevocationChecker,

		support.NewGinServer,
		support.NewGrpcServer,
		support.NewTelemetryManager,
		support.NewKafkaPublisher,
		infrastructure.NewInfra,
	),
	fx.Decorate(cache.NewCachingUserRepo),
	fx.Decorate(keycloak.NewCachingIntrospector),
	fx.Invoke(support.RegisterEventPublisher),
)

func NewApp(preConfig *pkgconfig.PreConfig, metadata *pkgconfig.Metadata) *fx.App {
	return fx.New(
		fx.Supply(preConfig, metadata),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
}
