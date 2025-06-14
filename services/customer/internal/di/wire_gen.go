// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
	config2 "github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/service"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/conductor"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/mysql"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/web"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeServer(env object.Env, metadata *config.Metadata) *infrastructure.Infra {
	configConfig := config2.NewConfig(env)
	serverConfig := config2.ProvideServerConfig(configConfig)
	appConfig := config2.ProvideAppConfig(configConfig)
	logHelper := logging.NewLogHelper(appConfig)
	grpcConfig := config2.ProvideGrpcConfig(configConfig)
	userGrpcClient := protostub.NewUserGrpcClient(grpcConfig)
	userMapper := mapper.NewUserMapper()
	customerUseCase := service.NewCustomerService(userGrpcClient, userMapper)
	customerAPI := web.NewCustomerController(customerUseCase, userMapper)
	orderGrpcClient := protostub.NewOrderGrpcClient(grpcConfig)
	orderUseCase := service.NewOrderService(orderGrpcClient)
	orderMapper := mapper.NewOrderMapper()
	orderAPI := web.NewOrderController(orderUseCase, orderMapper)
	paymentGrpcClient := protostub.NewPaymentGrpcClient(grpcConfig)
	paymentUseCase := service.NewPaymentService(paymentGrpcClient)
	paymentMapper := mapper.NewPaymentMapper()
	paymentAPI := web.NewPaymentController(paymentUseCase, paymentMapper)
	databaseConfig := config2.ProvideDatabaseConfig(configConfig)
	dataSource := database.NewDataSource(databaseConfig)
	pointRepo := mysql.NewPointRepository(dataSource)
	pointHandler := handler.NewPointHandler(pointRepo)
	pointUseCase := service.NewPointService(pointHandler)
	pointAPI := web.NewPointController(pointUseCase)
	wishlistRepo := mysql.NewWishlistRepository(dataSource)
	wishlistHandler := handler.NewWishlistHandler(wishlistRepo)
	wishlistUseCase := service.NewWishlistService(wishlistHandler)
	wishlistAPI := web.NewWishlistController(wishlistUseCase)
	ginServer := support.NewGinServer(serverConfig, logHelper, customerAPI, orderAPI, paymentAPI, pointAPI, wishlistAPI)
	conductorConfig := config2.ProvideConductorConfig(configConfig)
	workerClient := saga.NewWorkerClient(conductorConfig)
	customerWorker := conductor.NewCustomerWorker(conductorConfig)
	conductorWorker := support.NewConductorWorker(workerClient, customerWorker)
	telemetryConfig := config2.ProvideTelemetryConfig(configConfig)
	telemetryManager := support.NewTelemetryManager(metadata, telemetryConfig)
	infra := infrastructure.NewInfra(ginServer, conductorWorker, telemetryManager)
	return infra
}

func InitializeMigrationManager(env object.Env) *database.MigrationManager {
	configConfig := config2.NewConfig(env)
	databaseConfig := config2.ProvideDatabaseConfig(configConfig)
	migrationManager := database.NewMigrationManager(databaseConfig)
	return migrationManager
}

// wire.go:

var AppConfigProviderSet = wire.NewSet(config2.ProvideAppConfig)

var ServerConfigProviderSet = wire.NewSet(config2.ProvideServerConfig)

var DatabaseConfigProviderSet = wire.NewSet(config2.ProvideDatabaseConfig)

var ConductorConfigProviderSet = wire.NewSet(config2.ProvideConductorConfig)

var TelemetryConfigProviderSet = wire.NewSet(config2.ProvideTelemetryConfig)

var GrpcConfigProviderSet = wire.NewSet(config2.ProvideGrpcConfig)
