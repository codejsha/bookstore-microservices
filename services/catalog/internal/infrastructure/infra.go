package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support"
)

func NewInfra(
	ginServer *support.GinServer,
	grpcServer *support.GrpcServer,
	conductorWorker *support.ConductorWorker,
	telemetryManager *support.TelemetryManager,
) *Infra {
	return &Infra{
		ginServer:        ginServer,
		grpcServer:       grpcServer,
		conductorWorker:  conductorWorker,
		telemetryManager: telemetryManager,
	}
}

type Infra struct {
	ginServer        *support.GinServer
	grpcServer       *support.GrpcServer
	conductorWorker  *support.ConductorWorker
	telemetryManager *support.TelemetryManager
}

func (s *Infra) Run() {
	go s.telemetryManager.Run()
	go s.grpcServer.Run()
	go s.conductorWorker.Run()
	s.ginServer.Run()
}
