package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/support"
)

type Infra struct {
	ginServer        *support.GinServer
	grpcServer       *support.GrpcServer
	telemetryManager *support.TelemetryManager
}

func NewInfra(
	ginServer *support.GinServer,
	grpcServer *support.GrpcServer,
	telemetryManager *support.TelemetryManager,
) *Infra {
	return &Infra{
		ginServer:        ginServer,
		grpcServer:       grpcServer,
		telemetryManager: telemetryManager,
	}
}
