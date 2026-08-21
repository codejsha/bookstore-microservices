package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

func NewInfra(
	ginServer *support.GinServer,
	conductorWorker *support.ConductorWorker,
	telemetryManager *support.TelemetryManager,
) *Infra {
	return &Infra{
		ginServer:        ginServer,
		conductorWorker:  conductorWorker,
		telemetryManager: telemetryManager,
	}
}

type Infra struct {
	ginServer        *support.GinServer
	conductorWorker  *support.ConductorWorker
	telemetryManager *support.TelemetryManager
}

func (s *Infra) Run() {
	go s.telemetryManager.Run()
	go s.conductorWorker.Run()
	s.ginServer.Run()
}
