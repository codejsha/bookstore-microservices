package di

import (
	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/logging"
	"github.com/codejsha/shared-library-go/pkg/message"
	"github.com/codejsha/shared-library-go/pkg/rest/client"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/pgsql"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/restcontroller"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

var Module = fx.Module("customer",
	fx.Provide(
		config.NewConfig,
		config.ProvideAppConfig,
		config.ProvideServerConfig,
		config.ProvideDatabaseConfig,
		config.ProvideTelemetryConfig,
		config.ProvideGrpcConfig,
		config.ProvideKafkaConfig,

		pkgconfig.NewCloudConfigHelper,
		client.NewRestyClient,
		logging.NewLogHelper,

		restcontroller.NewCustomerController,
		restcontroller.NewOrderController,
		restcontroller.NewPaymentController,
		restcontroller.NewDeliveryController,
		restcontroller.NewPointController,
		restcontroller.NewReviewController,
		restcontroller.NewWishlistController,

		service.NewCustomerService,

		database.NewVaultAwareDataSource,
		support.NewReadinessDataSource,
		pgsql.NewPointRepository,
		pgsql.NewPointHistoryRepository,
		pgsql.NewReviewRepository,
		pgsql.NewWishlistRepository,

		protostub.NewOrderGrpcClient,
		protostub.NewPaymentGrpcClient,
		protostub.NewUserGrpcClient,
		protostub.NewDeliveryGrpcClient,
		protostub.NewUserClient,
		protostub.NewOrderClient,
		protostub.NewPaymentClient,
		protostub.NewDeliveryClient,

		support.NewGinServer,
		support.NewTelemetryManager,
		support.NewKafkaPublisher,
		infrastructure.NewInfra,
	),
	fx.Decorate(support.WithOutboundTimeout),
	fx.Invoke(support.ConfigureConnectionPool),
	fx.Invoke(support.RegisterEventPublisher),
	fx.Invoke(registerShutdownSequence),
)

func registerShutdownSequence(
	lc fx.Lifecycle,
	ginServer *support.GinServer,
	publisher *message.KafkaAsyncPublisher,
	readinessDataSource *support.ReadinessDataSource,
	telemetryManager *support.TelemetryManager,
	userClient *protostub.UserGrpcClient,
	orderClient *protostub.OrderGrpcClient,
	paymentClient *protostub.PaymentGrpcClient,
	deliveryClient *protostub.DeliveryGrpcClient,
) {
	support.RegisterShutdownSequence(lc, ginServer, publisher, readinessDataSource, telemetryManager,
		userClient, orderClient, paymentClient, deliveryClient)
}

func NewApp(preConfig *pkgconfig.PreConfig, metadata *pkgconfig.Metadata) *fx.App {
	return fx.New(
		fx.StopTimeout(support.ShutdownStopTimeout),
		fx.Supply(preConfig, metadata),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
}
