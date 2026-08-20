package infrastructure

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type Infra struct {
	ginServer        *support.GinServer
	telemetryManager *support.TelemetryManager
}

func NewInfra(
	ginServer *support.GinServer,
	telemetryManager *support.TelemetryManager,
) *Infra {
	return &Infra{
		ginServer:        ginServer,
		telemetryManager: telemetryManager,
	}
}
