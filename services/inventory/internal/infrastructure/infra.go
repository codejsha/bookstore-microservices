package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support"
)

type Infra struct {
	ginServer        *support.GinServer
	telemetryManager *support.TelemetryManager
	temporalWorker   *support.TemporalWorker
}

func NewInfra(
	ginServer *support.GinServer,
	telemetryManager *support.TelemetryManager,
	temporalWorker *support.TemporalWorker,
) *Infra {
	return &Infra{
		ginServer:        ginServer,
		telemetryManager: telemetryManager,
		temporalWorker:   temporalWorker,
	}
}
