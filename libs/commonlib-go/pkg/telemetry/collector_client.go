package telemetry

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func CreateCollectorClient(
	telemetryCfg *config.TelemetryConfig,
) (*grpc.ClientConn, error) {
	client, err := grpc.NewClient(telemetryCfg.Collector.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	return client, nil
}
