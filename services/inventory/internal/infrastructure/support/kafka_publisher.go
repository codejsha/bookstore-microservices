package support

import (
	"context"

	"go.uber.org/fx"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/restcontroller"
)

func NewKafkaPublisher(lc fx.Lifecycle, kafkaCfg *sharedconfig.KafkaConfig) (*message.KafkaAsyncPublisher, error) {
	pub, err := message.NewKafkaAsyncPublisher(kafkaCfg.Broker.BootstrapServers, message.NewKafkaAsyncProducerConfig())
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pub.Close()
		},
	})
	return pub, nil
}

func RegisterEventPublisher(pub *message.KafkaAsyncPublisher) {
	restcontroller.SetEventPublisher(pub)
}
