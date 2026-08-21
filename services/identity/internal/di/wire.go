//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/client"
	pkgconfig "github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/protosvc"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/web"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/support"
)

func InitializeServer(env object.Env, metadata *pkgconfig.Metadata) *infrastructure.Infra {
	wire.Build(
		config.NewConfig,
		AppConfigProviderSet,
		ServerConfigProviderSet,
		TelemetryConfigProviderSet,
		GrpcConfigProviderSet,

		logging.NewLogHelper,
		client.NewRestyClient,

		keycloak.NewAdminTokenHelper,
		keycloak.NewOpenIDConnectClient,
		keycloak.NewUsersClient,

		service.NewIdentityService,

		protosvc.NewUserGrpcServer,
		support.NewGrpcServer,

		web.NewIdentityController,
		support.NewGinServer,

		support.NewTelemetryManager,

		infrastructure.NewInfra,
	)
	return &infrastructure.Infra{}
}

var AppConfigProviderSet = wire.NewSet(config.ProvideAppConfig)
var ServerConfigProviderSet = wire.NewSet(config.ProvideServerConfig)
var TelemetryConfigProviderSet = wire.NewSet(config.ProvideTelemetryConfig)
var GrpcConfigProviderSet = wire.NewSet(config.ProvideGrpcConfig)
