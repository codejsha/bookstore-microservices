package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support"
)

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

type Infra struct {
	ginServer        *support.GinServer
	grpcServer       *support.GrpcServer
	telemetryManager *support.TelemetryManager
}

func (s *Infra) Run() {
	go s.telemetryManager.Run()
	go s.grpcServer.Run()
	s.ginServer.Run()
}
