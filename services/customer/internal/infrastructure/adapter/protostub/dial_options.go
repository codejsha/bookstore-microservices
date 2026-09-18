package protostub

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

const (
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second

	roundRobinServiceConfig = `{"loadBalancingConfig":[{"round_robin":{}}]}`
)

func keepaliveParameters() keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                keepaliveTime,
		Timeout:             keepaliveTimeout,
		PermitWithoutStream: true,
	}
}

func baseDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepaliveParameters()),
		grpc.WithDefaultServiceConfig(roundRobinServiceConfig),
	}
}

func dialOptions(telemetryManager *support.TelemetryManager) []grpc.DialOption {
	return append(baseDialOptions(), grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(telemetryManager.TraceProvider),
		otelgrpc.WithMeterProvider(telemetryManager.MeterProvider),
	)))
}
